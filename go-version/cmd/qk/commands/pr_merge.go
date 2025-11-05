package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Wangggym/quick-workflow/internal/git"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/spf13/cobra"
)

var prMergeCmd = &cobra.Command{
	Use:   "merge [pr-number]",
	Short: "Merge a PR and update Jira status",
	Long: `Merge a pull request and automatically:
  - Merge the PR on GitHub
  - Delete the remote branch
  - Delete the local branch
  - Update Jira status to Done/Merged`,
	Args: cobra.MaximumNArgs(1),
	Run:  runPRMerge,
}

func runPRMerge(cmd *cobra.Command, args []string) {
	// 检查是否在 Git 仓库中
	if !git.IsGitRepository() {
		ui.Error("Not a git repository")
		return
	}

	// 获取仓库信息
	remoteURL, err := git.GetRemoteURL()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get remote URL: %v", err))
		return
	}

	owner, repo, err := github.ParseRepositoryFromURL(remoteURL)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to parse repository: %v", err))
		return
	}

	// 获取 PR 号
	var prNumber int
	if len(args) > 0 {
		prNumber, err = strconv.Atoi(args[0])
		if err != nil {
			ui.Error(fmt.Sprintf("Invalid PR number: %s", args[0]))
			return
		}
	} else {
		prInput, err := ui.PromptInput("Enter PR number:", true)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to get input: %v", err))
			return
		}
		prNumber, err = strconv.Atoi(prInput)
		if err != nil {
			ui.Error(fmt.Sprintf("Invalid PR number: %s", prInput))
			return
		}
	}

	// 创建 GitHub 客户端
	ghClient, err := github.NewClient()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create GitHub client: %v", err))
		return
	}

	// 获取 PR 信息
	ui.Info(fmt.Sprintf("Fetching PR #%d...", prNumber))
	pr, err := ghClient.GetPullRequest(owner, repo, prNumber)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get PR: %v", err))
		return
	}

	ui.Info(fmt.Sprintf("PR: %s", pr.Title))
	ui.Info(fmt.Sprintf("Branch: %s -> %s", pr.Head, pr.Base))

	// 确认合并
	confirm, err := ui.PromptConfirm(fmt.Sprintf("Merge PR #%d?", prNumber), true)
	if err != nil || !confirm {
		ui.Info("Merge cancelled")
		return
	}

	// 合并 PR
	ui.Info("Merging pull request...")
	if err := ghClient.MergePullRequest(owner, repo, prNumber, pr.Title); err != nil {
		ui.Error(fmt.Sprintf("Failed to merge PR: %v", err))
		return
	}
	ui.Success("Pull request merged!")

	// 删除远程分支
	deleteBranch, err := ui.PromptConfirm("Delete remote branch?", true)
	if err == nil && deleteBranch {
		ui.Info(fmt.Sprintf("Deleting remote branch %s...", pr.Head))
		if err := git.DeleteRemoteBranch(pr.Head); err != nil {
			ui.Warning(fmt.Sprintf("Failed to delete remote branch: %v", err))
		} else {
			ui.Success("Remote branch deleted")
		}
	}

	// 切换到主分支
	currentBranch, err := git.GetCurrentBranch()
	if err == nil && currentBranch == pr.Head {
		ui.Info("Switching to main branch...")
		if err := git.CreateBranch("main"); err != nil {
			// 可能已经在 main 上，尝试 checkout
			ui.Warning("Could not switch to main, you may need to do this manually")
		}
	}

	// 删除本地分支
	deleteLocal, err := ui.PromptConfirm("Delete local branch?", true)
	if err == nil && deleteLocal {
		ui.Info(fmt.Sprintf("Deleting local branch %s...", pr.Head))
		if err := git.DeleteBranch(pr.Head); err != nil {
			ui.Warning(fmt.Sprintf("Failed to delete local branch: %v", err))
		} else {
			ui.Success("Local branch deleted")
		}
	}

	// 从标题中提取 Jira ticket
	jiraTicket := extractJiraTicket(pr.Title)
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		ui.Info(fmt.Sprintf("Found Jira ticket: %s", jiraTicket))

		updateJira, err := ui.PromptConfirm("Update Jira status?", true)
		if err == nil && updateJira {
			jiraClient, err := jira.NewClient()
			if err != nil {
				ui.Warning(fmt.Sprintf("Failed to create Jira client: %v", err))
			} else {
				projectKey := jira.ExtractProjectKey(jiraTicket)
				statuses, err := jiraClient.GetProjectStatuses(projectKey)
				if err != nil {
					ui.Warning(fmt.Sprintf("Failed to get statuses: %v", err))
				} else {
					// 默认选择 "Done" 或类似的状态
					defaultStatus := findDefaultMergedStatus(statuses)
					newStatus, err := ui.PromptSelect("Select new status:", statuses)
					if err == nil {
						if newStatus == "" {
							newStatus = defaultStatus
						}
						if err := jiraClient.UpdateStatus(jiraTicket, newStatus); err != nil {
							ui.Warning(fmt.Sprintf("Failed to update status: %v", err))
						} else {
							ui.Success(fmt.Sprintf("Updated Jira status to: %s", newStatus))
						}
					}
				}

				// 添加合并评论
				comment := fmt.Sprintf("PR #%d merged: %s", prNumber, pr.HTMLURL)
				if err := jiraClient.AddComment(jiraTicket, comment); err != nil {
					ui.Warning(fmt.Sprintf("Failed to add comment: %v", err))
				}
			}
		}
	}

	fmt.Println()
	ui.Success("All done! 🎉")
}

func extractJiraTicket(title string) string {
	// 尝试从标题中提取 Jira ticket，格式通常是 "PROJ-123: Title"
	parts := strings.Split(title, ":")
	if len(parts) > 0 {
		candidate := strings.TrimSpace(parts[0])
		if jira.ValidateIssueKey(candidate) {
			return candidate
		}
	}
	return ""
}

func findDefaultMergedStatus(statuses []string) string {
	// 查找 "Done" 或类似的状态
	lowerStatuses := make(map[string]string)
	for _, s := range statuses {
		lowerStatuses[strings.ToLower(s)] = s
	}

	preferredStatuses := []string{"done", "merged", "closed", "resolved"}
	for _, preferred := range preferredStatuses {
		if status, ok := lowerStatuses[preferred]; ok {
			return status
		}
	}

	if len(statuses) > 0 {
		return statuses[len(statuses)-1] // 返回最后一个
	}
	return ""
}

