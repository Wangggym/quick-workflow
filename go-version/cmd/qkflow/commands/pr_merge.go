package commands

import (
	"fmt"
	"os"
	"os/exec"
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

	// 创建 GitHub 客户端
	ghClient, err := github.NewClient()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create GitHub client: %v", err))
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
		// 尝试获取当前分支的 PR
		currentBranch, err := git.GetCurrentBranch()
		if err == nil && currentBranch != "" {
			ui.Info(fmt.Sprintf("Checking for PR from current branch: %s", currentBranch))
			
			// 先尝试 open 状态的 PR
			prs, err := ghClient.ListPullRequests(owner, repo, "open")
			if err == nil {
				for _, pr := range prs {
					if pr.Head == currentBranch {
						prNumber = pr.Number
						ui.Success(fmt.Sprintf("Found PR #%d: %s", pr.Number, pr.Title))
						break
					}
				}
			}
			
			// 如果没找到，尝试所有状态的 PR
			if prNumber == 0 {
				allPRs, err := ghClient.ListPullRequests(owner, repo, "all")
				if err == nil {
					for _, pr := range allPRs {
						if pr.Head == currentBranch {
							prNumber = pr.Number
							ui.Success(fmt.Sprintf("Found PR #%d (%s): %s", pr.Number, pr.State, pr.Title))
							break
						}
					}
				}
			}
			
			// 如果还是没找到，提示用户
			if prNumber == 0 {
				ui.Warning(fmt.Sprintf("No PR found for branch: %s", currentBranch))
				ui.Info("This branch may not have a PR yet. Please create one first with:")
				ui.Info("  qkg pr create")
				fmt.Println()
				
				// 询问用户是否手动输入 PR 号
				manually, err := ui.PromptConfirm("Do you want to manually enter a PR number or select from list?", true)
				if err != nil || !manually {
					ui.Info("Merge cancelled")
					return
				}
			}
		}

		// 如果用户选择手动输入或没有当前分支
		if prNumber == 0 {
			prs, err := ghClient.ListPullRequests(owner, repo, "open")
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to list PRs: %v", err))
				return
			}

			if len(prs) == 0 {
				ui.Error("No open pull requests found")
				return
			}

			// 构建选择列表
			prOptions := make([]string, len(prs))
			for i, pr := range prs {
				prOptions[i] = fmt.Sprintf("#%d - %s", pr.Number, pr.Title)
			}

			selected, err := ui.PromptSelect("Select a PR to merge:", prOptions)
			if err != nil {
				if err.Error() == "interrupt" {
					ui.Warning("Operation cancelled by user")
					os.Exit(0)
				}
				ui.Error(fmt.Sprintf("Failed to select PR: %v", err))
				return
			}

			// 从选择中提取 PR 号
			var selectedPR *github.PullRequest
			for i, option := range prOptions {
				if option == selected {
					selectedPR = &prs[i]
					break
				}
			}

			if selectedPR != nil {
				prNumber = selectedPR.Number
			} else {
				ui.Error("Failed to find selected PR")
				return
			}
		}
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
	ui.Info(fmt.Sprintf("State: %s", pr.State))

	// 检查 PR 状态
	alreadyMerged := false
	if pr.State == "closed" {
		ui.Warning("This PR is already closed")
		// 检查是否是已合并
		alreadyMerged = true
	} else {
		// 合并 PR
		ui.Info(fmt.Sprintf("Merging PR #%d...", prNumber))
		if err := ghClient.MergePullRequest(owner, repo, prNumber, pr.Title); err != nil {
			ui.Error(fmt.Sprintf("Failed to merge PR: %v", err))
			return
		}
		ui.Success("Pull request merged!")
	}

	// 删除远程分支（如果还存在）
	if !alreadyMerged {
		ui.Info(fmt.Sprintf("Deleting remote branch %s...", pr.Head))
		if err := git.DeleteRemoteBranch(pr.Head); err != nil {
			ui.Warning(fmt.Sprintf("Failed to delete remote branch: %v (may already be deleted)", err))
		} else {
			ui.Success("Remote branch deleted")
		}
	} else {
		ui.Info("Skipping remote branch deletion (PR already merged)")
	}

	// 切换到主分支并删除本地分支
	currentBranch, err := git.GetCurrentBranch()
	if err == nil && currentBranch == pr.Head {
		// 获取默认分支
		defaultBranch, err := git.GetDefaultBranch()
		if err != nil {
			defaultBranch = "master"
		}
		
		ui.Info(fmt.Sprintf("Switching to %s branch...", defaultBranch))
		// 使用 checkout 而不是 create
		cmd := exec.Command("git", "checkout", defaultBranch)
		if err := cmd.Run(); err != nil {
			ui.Warning(fmt.Sprintf("Could not switch to %s, you may need to do this manually", defaultBranch))
		} else {
			// 切换成功后，拉取最新代码
			ui.Info(fmt.Sprintf("Pulling latest changes from %s...", defaultBranch))
			pullCmd := exec.Command("git", "pull")
			if err := pullCmd.Run(); err != nil {
				ui.Warning("Failed to pull latest changes, you may need to run 'git pull' manually")
			} else {
				ui.Success("Updated to latest changes")
			}
		}
		
		// 删除本地分支
		ui.Info(fmt.Sprintf("Deleting local branch %s...", pr.Head))
		if err := git.DeleteBranch(pr.Head); err != nil {
			ui.Warning(fmt.Sprintf("Failed to delete local branch: %v", err))
		} else {
			ui.Success("Local branch deleted")
		}
	}

	// 从标题中提取 Jira ticket 并自动更新
	jiraTicket := extractJiraTicket(pr.Title)
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		ui.Info(fmt.Sprintf("Found Jira ticket: %s", jiraTicket))

		jiraClient, err := jira.NewClient()
		if err != nil {
			ui.Warning(fmt.Sprintf("Failed to create Jira client: %v", err))
		} else {
			// 使用缓存的状态
			projectKey := jira.ExtractProjectKey(jiraTicket)
			
			statusCache, err := jira.NewStatusCache()
			if err != nil {
				ui.Warning(fmt.Sprintf("Failed to create status cache: %v", err))
			} else {
				mapping, err := statusCache.GetProjectStatus(projectKey)
				if err != nil {
					ui.Warning(fmt.Sprintf("Failed to get cached status: %v", err))
				} else if mapping != nil && mapping.PRMergedStatus != "" {
					// 使用缓存的 merged 状态
					ui.Info(fmt.Sprintf("Updating Jira status to: %s", mapping.PRMergedStatus))
					if err := jiraClient.UpdateStatus(jiraTicket, mapping.PRMergedStatus); err != nil {
						ui.Warning(fmt.Sprintf("Failed to update status: %v", err))
					} else {
						ui.Success(fmt.Sprintf("Updated Jira status to: %s", mapping.PRMergedStatus))
					}
				} else {
					// 没有缓存，使用默认逻辑
					statuses, err := jiraClient.GetProjectStatuses(projectKey)
					if err != nil {
						ui.Warning(fmt.Sprintf("Failed to get statuses: %v", err))
					} else {
						defaultStatus := findDefaultMergedStatus(statuses)
						if defaultStatus != "" {
							ui.Info(fmt.Sprintf("Updating Jira status to: %s", defaultStatus))
							if err := jiraClient.UpdateStatus(jiraTicket, defaultStatus); err != nil {
								ui.Warning(fmt.Sprintf("Failed to update status: %v", err))
							} else {
								ui.Success(fmt.Sprintf("Updated Jira status to: %s", defaultStatus))
							}
						}
					}
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

