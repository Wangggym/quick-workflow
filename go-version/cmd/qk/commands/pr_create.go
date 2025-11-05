package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Wangggym/quick-workflow/internal/git"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/Wangggym/quick-workflow/pkg/config"
	"github.com/spf13/cobra"
)

var prCreateCmd = &cobra.Command{
	Use:   "create [jira-ticket]",
	Short: "Create a PR and update Jira status",
	Long: `Create a new pull request and automatically:
  - Create a git branch
  - Commit staged changes
  - Push to remote
  - Create a GitHub PR
  - Add PR link to Jira
  - Update Jira status`,
	Args: cobra.MaximumNArgs(1),
	Run:  runPRCreate,
}

func runPRCreate(cmd *cobra.Command, args []string) {
	// 检查是否在 Git 仓库中
	if !git.IsGitRepository() {
		ui.Error("Not a git repository")
		return
	}

	// 检查是否有未提交的更改
	hasChanges, err := git.HasUncommittedChanges()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check git status: %v", err))
		return
	}
	if !hasChanges {
		ui.Error("No changes to commit. Please stage your changes first with 'git add'")
		return
	}

	// 获取或输入 Jira ticket
	var jiraTicket string
	if len(args) > 0 {
		jiraTicket = args[0]
	} else {
		jiraTicket, err = ui.PromptInput("Jira ticket (optional, press Enter to skip):", false)
		if err != nil {
			// 用户取消操作
			if err.Error() == "interrupt" {
				ui.Warning("Operation cancelled by user")
				os.Exit(0)
			}
			ui.Error(fmt.Sprintf("Failed to get input: %v", err))
			return
		}
	}

	// 如果有 Jira ticket，获取信息
	var jiraIssue *jira.Issue
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		jiraClient, err := jira.NewClient()
		if err != nil {
			ui.Warning(fmt.Sprintf("Failed to create Jira client: %v", err))
		} else {
			jiraIssue, err = jiraClient.GetIssue(jiraTicket)
			if err != nil {
				ui.Warning(fmt.Sprintf("Failed to get Jira issue: %v", err))
			} else {
				ui.Info(fmt.Sprintf("Found Jira issue: %s", jiraIssue.Summary))
			}
		}
	}

	// 获取 PR 标题
	var title string
	if jiraIssue != nil {
		title = jiraIssue.Summary
		ui.Info(fmt.Sprintf("Using Jira summary as title: %s", title))
		confirm, err := ui.PromptConfirm("Use this title?", true)
		if err != nil {
			if err.Error() == "interrupt" {
				ui.Warning("Operation cancelled by user")
				os.Exit(0)
			}
			ui.Error(fmt.Sprintf("Failed to get confirmation: %v", err))
			return
		}
		if !confirm {
			title, err = ui.PromptInput("Enter PR title:", true)
			if err != nil {
				if err.Error() == "interrupt" {
					ui.Warning("Operation cancelled by user")
					os.Exit(0)
				}
				ui.Error(fmt.Sprintf("Failed to get title: %v", err))
				return
			}
		}
	} else {
		title, err = ui.PromptInput("Enter PR title:", true)
		if err != nil {
			if err.Error() == "interrupt" {
				ui.Warning("Operation cancelled by user")
				os.Exit(0)
			}
			ui.Error(fmt.Sprintf("Failed to get title: %v", err))
			return
		}
	}

	// 获取描述
	description, err := ui.PromptInput("Enter PR description (optional):", false)
	if err != nil {
		description = ""
	}

	// 选择变更类型
	prTypes := ui.PRTypeOptions()
	selectedTypes, err := ui.PromptMultiSelect("Select type(s) of changes:", prTypes)
	if err != nil {
		if err.Error() == "interrupt" {
			ui.Warning("Operation cancelled by user")
			os.Exit(0)
		}
		ui.Warning("No types selected, continuing...")
		selectedTypes = []string{}
	}

	// 构建 PR body
	prBody := buildPRBody(description, selectedTypes, jiraTicket)

	// 创建分支名
	branchName := buildBranchName(jiraTicket, title)
	cfg := config.Get()
	if cfg.BranchPrefix != "" {
		branchName = cfg.BranchPrefix + "/" + branchName
	}

	ui.Info(fmt.Sprintf("Creating branch: %s", branchName))

	// 创建分支
	if err := git.CreateBranch(branchName); err != nil {
		ui.Error(fmt.Sprintf("Failed to create branch: %v", err))
		return
	}

	// Stage 所有更改
	ui.Info("Staging changes...")
	if err := git.AddAll(); err != nil {
		ui.Error(fmt.Sprintf("Failed to stage changes: %v", err))
		return
	}

	// 提交更改
	commitMessage := title
	if jiraTicket != "" {
		commitMessage = fmt.Sprintf("%s: %s", jiraTicket, title)
	} else {
		// 无 Jira ticket 时，添加 # 前缀
		commitMessage = fmt.Sprintf("# %s", title)
	}
	
	ui.Info("Committing changes...")
	if err := git.Commit(commitMessage); err != nil {
		ui.Error(fmt.Sprintf("Failed to commit: %v", err))
		return
	}

	// 推送分支
	ui.Info("Pushing branch to remote...")
	if err := git.Push(branchName); err != nil {
		ui.Error(fmt.Sprintf("Failed to push: %v", err))
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

	// 获取默认分支
	defaultBranch, err := git.GetDefaultBranch()
	if err != nil {
		ui.Warning(fmt.Sprintf("Failed to detect default branch, using 'main': %v", err))
		defaultBranch = "main"
	}
	ui.Info(fmt.Sprintf("Using base branch: %s", defaultBranch))

	// 创建 PR
	ui.Info("Creating pull request...")
	ghClient, err := github.NewClient()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create GitHub client: %v", err))
		return
	}

	pr, err := ghClient.CreatePullRequest(github.CreatePullRequestInput{
		Owner: owner,
		Repo:  repo,
		Title: commitMessage,
		Body:  prBody,
		Head:  branchName,
		Base:  defaultBranch,
	})
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create PR: %v", err))
		return
	}

	ui.Success(fmt.Sprintf("Pull request created: %s", pr.HTMLURL))

	// 更新 Jira
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		jiraClient, err := jira.NewClient()
		if err != nil {
			ui.Warning(fmt.Sprintf("Failed to create Jira client: %v", err))
		} else {
			// 添加 PR 链接
			ui.Info("Adding PR link to Jira...")
			if err := jiraClient.AddPRLink(jiraTicket, pr.HTMLURL); err != nil {
				ui.Warning(fmt.Sprintf("Failed to add PR link to Jira: %v", err))
			} else {
				ui.Success("Added PR link to Jira")
			}

			// 更新状态
			updateStatus, err := ui.PromptConfirm("Do you want to update Jira status?", true)
			if err == nil && updateStatus {
				projectKey := jira.ExtractProjectKey(jiraTicket)
				statuses, err := jiraClient.GetProjectStatuses(projectKey)
				if err != nil {
					ui.Warning(fmt.Sprintf("Failed to get statuses: %v", err))
				} else {
					newStatus, err := ui.PromptSelect("Select new status:", statuses)
					if err == nil {
						if err := jiraClient.UpdateStatus(jiraTicket, newStatus); err != nil {
							ui.Warning(fmt.Sprintf("Failed to update status: %v", err))
						} else {
							ui.Success(fmt.Sprintf("Updated Jira status to: %s", newStatus))
						}
					}
				}
			}
		}
	}

	// 复制 URL 到剪贴板
	copyToClipboard(pr.HTMLURL)
	
	// 打开浏览器
	openBrowser(pr.HTMLURL)
	
	fmt.Println()
	ui.Success("All done! 🎉")
}

func buildBranchName(jiraTicket, title string) string {
	sanitized := git.SanitizeBranchName(title)
	if jiraTicket != "" {
		return fmt.Sprintf("%s--%s", jiraTicket, sanitized)
	}
	return sanitized
}

func buildPRBody(description string, types []string, jiraTicket string) string {
	var body strings.Builder

	if description != "" {
		body.WriteString(description)
		body.WriteString("\n\n")
	}

	if len(types) > 0 {
		body.WriteString("## Types of changes\n\n")
		for _, t := range types {
			body.WriteString(fmt.Sprintf("- [x] %s\n", t))
		}
		body.WriteString("\n")
	}

	if jiraTicket != "" {
		body.WriteString(fmt.Sprintf("## Related Issues\n\n- %s\n", jiraTicket))
	}

	return body.String()
}

func copyToClipboard(text string) {
	// macOS only for now
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		ui.Warning("Failed to copy to clipboard")
	} else {
		ui.Success(fmt.Sprintf("Successfully copied %s to clipboard", text))
	}
}

func openBrowser(url string) {
	// macOS
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		ui.Warning(fmt.Sprintf("Failed to open browser: %v", err))
	}
}

