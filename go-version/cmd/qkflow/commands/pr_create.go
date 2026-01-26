package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Wangggym/quick-workflow/internal/ai"
	"github.com/Wangggym/quick-workflow/internal/git"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/Wangggym/quick-workflow/internal/watcher"
	"github.com/Wangggym/quick-workflow/pkg/config"
	"github.com/spf13/cobra"
)

var (
	prDesc   string
	prTypes  []string
	noTicket bool
	prTitle  string
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

func init() {
	prCreateCmd.Flags().StringVar(&prDesc, "pr-desc", "", "PR description body content")
	prCreateCmd.Flags().StringSliceVar(&prTypes, "types", []string{}, "Change types (e.g., feat,fix,docs)")
	prCreateCmd.Flags().BoolVar(&noTicket, "no-ticket", false, "Skip Jira ticket (proceed without ticket)")
	prCreateCmd.Flags().StringVar(&prTitle, "title", "", "PR title (if not provided, will be generated from description)")
}

func runPRCreate(cmd *cobra.Command, args []string) {
	// 检查是否在 Git 仓库中
	if !git.IsGitRepository() {
		ui.Error("Not a git repository")
		return
	}

	// 记录原始分支，以便失败时回退
	originalBranch, err := git.GetCurrentBranch()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get current branch: %v", err))
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
	if noTicket {
		// 明确指定跳过 ticket
		jiraTicket = ""
	} else if len(args) > 0 {
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

	// 显示 Jira 信息
	if jiraIssue != nil {
		ui.Info(fmt.Sprintf("Jira issue: %s", jiraIssue.Summary))
	}

	// 选择变更类型
	var selectedTypes []string
	if len(prTypes) > 0 {
		// 使用命令行提供的类型
		selectedTypes = prTypes
		ui.Info(fmt.Sprintf("Using provided types: %v", prTypes))
	} else {
		// 交互式选择
		prTypeOptions := ui.PRTypeOptions()
		selectedTypes, err = ui.PromptMultiSelect("Select type(s) of changes:", prTypeOptions)
		if err != nil {
			if err.Error() == "interrupt" {
				ui.Warning("Operation cancelled by user")
				os.Exit(0)
			}
			ui.Warning("No types selected, continuing...")
			selectedTypes = []string{}
		}
	}

	// 生成 PR 标题
	var title string
	if jiraIssue != nil {
		// 提取主要类型（第一个选择的类型）
		prType := ""
		if len(selectedTypes) > 0 {
			prType = ui.ExtractPRType(selectedTypes[0])
		}

		// 使用 AI 生成简洁的 PR 标题
		aiClient, err := ai.NewClient()
		if err == nil && prType != "" {
			ui.Info("Generating PR title with AI...")
			title, err = aiClient.GeneratePRTitle(jiraIssue.Summary, prType, "")
			if err != nil {
				ui.Warning(fmt.Sprintf("AI generation failed: %v", err))
				// 回退到简单格式
				title = generateSimpleTitle(jiraIssue.Summary, prType, "")
			} else {
				// 确保标题有类型前缀且格式正确
				title = ensureTitlePrefix(title, prType)
				ui.Success(fmt.Sprintf("Generated title: %s", title))
			}
		} else {
			// 没有 AI 或没有类型，使用简单生成
			title = generateSimpleTitle(jiraIssue.Summary, prType, "")
			ui.Success(fmt.Sprintf("Generated title: %s", title))
		}
	} else {
		// 没有 Jira
		if prTitle != "" {
			// 使用提供的标题，确保有类型前缀
			if len(selectedTypes) > 0 {
				prType := ui.ExtractPRType(selectedTypes[0])
				title = ensureTitlePrefix(prTitle, prType)
			} else {
				title = prTitle
			}
			ui.Success(fmt.Sprintf("Using provided title: %s", title))
		} else if prDesc != "" {
			// 从描述中提取标题
			lines := strings.Split(strings.TrimSpace(prDesc), "\n")
			var titleText string

			// 查找 "## Summary" 后面的第一个非空行
			for i, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "## Summary") {
					// 跳过 Summary 行和后续的空行，找到第一个非空行
					for j := i + 1; j < len(lines); j++ {
						nextLine := strings.TrimSpace(lines[j])
						if nextLine != "" && !strings.HasPrefix(nextLine, "#") {
							titleText = nextLine
							break
						}
					}
					if titleText != "" {
						break
					}
				} else if !strings.HasPrefix(line, "#") && line != "" && titleText == "" {
					// 使用第一个非标题行
					titleText = line
					break
				}
			}

			// 如果没有找到，使用第一行
			if titleText == "" && len(lines) > 0 {
				titleText = strings.TrimSpace(lines[0])
				titleText = strings.TrimPrefix(titleText, "# ")
				titleText = strings.TrimPrefix(titleText, "## ")
			}

			// 确保标题有类型前缀
			if len(selectedTypes) > 0 {
				prType := ui.ExtractPRType(selectedTypes[0])
				// 如果标题已经有类型前缀，不重复添加
				if !strings.HasPrefix(strings.ToLower(titleText), strings.ToLower(prType)+":") {
					title = fmt.Sprintf("%s: %s", capitalizeFirst(prType), truncateString(titleText, 80))
				} else {
					title = truncateString(titleText, 100)
				}
			} else {
				title = truncateString(titleText, 100)
			}
			ui.Success(fmt.Sprintf("Generated title: %s", title))
		} else {
			// 没有描述，手动输入
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
	}

	// 构建 PR body（包含描述）
	prBody := buildPRBody(selectedTypes, jiraTicket, prDesc)

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
		rollbackBranch(originalBranch, branchName)
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
		// 重试一次
		ui.Warning(fmt.Sprintf("Failed to create PR: %v", err))
		ui.Info("Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)

		ui.Info("Retrying to create pull request...")
		pr, err = ghClient.CreatePullRequest(github.CreatePullRequestInput{
			Owner: owner,
			Repo:  repo,
			Title: commitMessage,
			Body:  prBody,
			Head:  branchName,
			Base:  defaultBranch,
		})
		if err != nil {
			ui.Error(fmt.Sprintf("Retry failed: %v", err))
			ui.Info("Rolling back changes...")
			rollbackBranch(originalBranch, branchName)
			return
		}
	}

	ui.Success(fmt.Sprintf("Pull request created: %s", pr.HTMLURL))

	// 如果有 Jira ticket，添加 PR 链接到 Jira（描述已经在 PR body 中了）
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		jiraClient, err := jira.NewClient()
		if err == nil {
			ui.Info("Adding PR link to Jira...")
			jiraComment := fmt.Sprintf("*PR Created:*\n\n[View PR|%s]", pr.HTMLURL)
			if err := jiraClient.AddComment(jiraTicket, jiraComment); err != nil {
				ui.Warning(fmt.Sprintf("Failed to add comment to Jira: %v", err))
			} else {
				ui.Success("PR link added to Jira")
			}
		}
	}

	// 更新 Jira
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		jiraClient, err := jira.NewClient()
		if err != nil {
			ui.Warning(fmt.Sprintf("Failed to create Jira client: %v", err))
		} else {
			// 分配给当前用户
			ui.Info("Assigning Jira ticket to you...")
			if err := jiraClient.AssignToMe(jiraTicket); err != nil {
				ui.Warning(fmt.Sprintf("Failed to assign ticket: %v", err))
			} else {
				ui.Success("Assigned Jira ticket to you")
			}

			// 添加 PR 链接
			ui.Info("Adding PR link to Jira...")
			if err := jiraClient.AddPRLink(jiraTicket, pr.HTMLURL); err != nil {
				ui.Warning(fmt.Sprintf("Failed to add PR link to Jira: %v", err))
			} else {
				ui.Success("Added PR link to Jira")
			}

			// 更新状态
			projectKey := jira.ExtractProjectKey(jiraTicket)

			// 检查状态缓存
			statusCache, err := jira.NewStatusCache()
			if err != nil {
				ui.Warning(fmt.Sprintf("Failed to create status cache: %v", err))
			} else {
				mapping, err := statusCache.GetProjectStatus(projectKey)
				if err != nil {
					ui.Warning(fmt.Sprintf("Failed to get cached status: %v", err))
				} else if mapping == nil {
					// 第一次使用，需要配置状态映射
					// 检查是否有提供 --types 或 --pr-desc（表示可能是非交互模式）
					if len(prTypes) > 0 || prDesc != "" {
						ui.Error(fmt.Sprintf("❌ No status mapping found for project %s", projectKey))
						ui.Info("💡 Please run 'qkflow pr create' interactively first to configure status mappings")
						ui.Info("   Then you can use --types and --pr-desc flags for automation")
					} else {
						// 交互式配置状态映射
						ui.Info(fmt.Sprintf("First time using project %s, please configure status mappings", projectKey))
						mapping, err = setupProjectStatusMapping(jiraClient, projectKey)
						if err != nil {
							ui.Warning(fmt.Sprintf("Failed to setup status mapping: %v", err))
						} else if mapping != nil {
							// 保存配置
							if err := statusCache.SaveProjectStatus(mapping); err != nil {
								ui.Warning(fmt.Sprintf("Failed to save status mapping: %v", err))
							} else {
								ui.Success("Status mapping saved!")
							}
						}
					}
				}

				// 使用缓存的状态更新
				if mapping != nil && mapping.PRCreatedStatus != "" {
					ui.Info(fmt.Sprintf("Updating Jira status to: %s", mapping.PRCreatedStatus))
					if err := jiraClient.UpdateStatus(jiraTicket, mapping.PRCreatedStatus); err != nil {
						ui.Warning(fmt.Sprintf("Failed to update status: %v", err))
					} else {
						ui.Success(fmt.Sprintf("Updated Jira status to: %s", mapping.PRCreatedStatus))
					}
				}
			}
		}
	}

	// 添加到 watching list
	watchingList, err := watcher.NewWatchingList()
	if err != nil {
		ui.Warning(fmt.Sprintf("Failed to load watching list: %v", err))
	} else {
		// Extract Jira tickets
		jiraTickets := make([]string, 0)
		if jiraTicket != "" {
			jiraTickets = append(jiraTickets, jiraTicket)
		}

		watchingPR := watcher.WatchingPR{
			PRNumber:    pr.Number,
			Owner:       owner,
			Repo:        repo,
			Branch:      branchName,
			Title:       commitMessage,
			PRURL:       pr.HTMLURL,
			JiraTickets: jiraTickets,
		}

		if err := watchingList.Add(watchingPR); err != nil {
			ui.Warning(fmt.Sprintf("Failed to add PR to watching list: %v", err))
		} else {
			ui.Info("✅ Added PR to watching list for auto Jira updates")
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

func buildPRBody(types []string, jiraTicket string, prDesc string) string {
	var body strings.Builder

	body.WriteString("# PR Ready\n\n")

	// 检查 prDesc 中是否已经包含 "Types of changes" 部分（不区分大小写）
	prDescLower := strings.ToLower(prDesc)
	hasTypesInDesc := strings.Contains(prDescLower, "types of changes") || 
		strings.Contains(prDescLower, "## types of changes")

	if len(types) > 0 && !hasTypesInDesc {
		// prDesc 中没有 Types of changes，自动添加
		body.WriteString("## Types of changes\n\n")
		for _, t := range types {
			body.WriteString(fmt.Sprintf("- [x] %s\n", t))
		}
		body.WriteString("\n")
	}

	if jiraTicket != "" {
		cfg := config.Get()
		jiraURL := fmt.Sprintf("%s/browse/%s", cfg.JiraServiceAddress, jiraTicket)
		body.WriteString(fmt.Sprintf("#### Jira Link:\n\n%s\n\n", jiraURL))
	}

	// 如果提供了描述，添加到 body 中
	if prDesc != "" {
		body.WriteString("---\n\n")
		body.WriteString(prDesc)
		body.WriteString("\n")
	}

	return body.String()
}

func setupProjectStatusMapping(client *jira.Client, projectKey string) (*jira.StatusMapping, error) {
	statuses, err := client.GetProjectStatuses(projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get project statuses: %w", err)
	}

	ui.Info("Select status when PR is created/in progress:")
	createdStatus, err := ui.PromptSelect("Status for PR created:", statuses)
	if err != nil {
		if err.Error() == "interrupt" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select created status: %w", err)
	}

	ui.Info("Select status when PR is merged/done:")
	mergedStatus, err := ui.PromptSelect("Status for PR merged:", statuses)
	if err != nil {
		if err.Error() == "interrupt" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select merged status: %w", err)
	}

	return &jira.StatusMapping{
		ProjectKey:      projectKey,
		PRCreatedStatus: createdStatus,
		PRMergedStatus:  mergedStatus,
	}, nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func ensureTitlePrefix(title, prType string) string {
	if prType == "" {
		return title
	}

	// 检查标题是否已经有类型前缀
	prTypeLower := strings.ToLower(prType)
	titleLower := strings.ToLower(title)

	// 检查各种可能的格式
	if strings.HasPrefix(titleLower, prTypeLower+":") ||
		strings.HasPrefix(titleLower, capitalizeFirst(prType)+":") {
		// 已经有前缀，确保格式正确（首字母大写）
		parts := strings.SplitN(title, ":", 2)
		if len(parts) == 2 {
			return fmt.Sprintf("%s:%s", capitalizeFirst(prType), parts[1])
		}
		return title
	}

	// 没有前缀，添加
	return fmt.Sprintf("%s: %s", capitalizeFirst(prType), title)
}

func generateSimpleTitle(jiraSummary, prType, description string) string {
	// 如果有简短描述，使用描述
	if description != "" {
		if prType != "" {
			return fmt.Sprintf("%s: %s", capitalizeFirst(prType), description)
		}
		return description
	}

	// 否则使用 Jira 标题的前 50 个字符
	summary := jiraSummary
	if len(summary) > 50 {
		summary = summary[:50] + "..."
	}

	if prType != "" {
		return fmt.Sprintf("%s: %s", capitalizeFirst(prType), summary)
	}
	return summary
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

// rollbackBranch rolls back branch creation when PR creation fails
// It preserves changes, deletes the remote/local branch, and switches back to original branch
func rollbackBranch(originalBranch, newBranch string) {
	// 1. 撤销 commit 但保留改动在暂存区
	ui.Info("Undoing commit to preserve changes...")
	resetCmd := exec.Command("git", "reset", "--soft", "HEAD~1")
	if err := resetCmd.Run(); err != nil {
		ui.Warning(fmt.Sprintf("Failed to reset commit: %v", err))
	}

	// 2. 暂存改动
	ui.Info("Stashing changes...")
	stashCmd := exec.Command("git", "stash", "--include-untracked")
	if err := stashCmd.Run(); err != nil {
		ui.Warning(fmt.Sprintf("Failed to stash changes: %v", err))
	}

	// 3. 删除远程分支
	ui.Info("Deleting remote branch...")
	if err := git.DeleteRemoteBranch(newBranch); err != nil {
		ui.Warning(fmt.Sprintf("Failed to delete remote branch: %v", err))
	} else {
		ui.Success("Remote branch deleted")
	}

	// 4. 切换回原分支
	ui.Info(fmt.Sprintf("Switching back to branch: %s", originalBranch))
	if err := git.CheckoutBranch(originalBranch); err != nil {
		ui.Error(fmt.Sprintf("Failed to checkout original branch: %v", err))
		ui.Warning("Your changes are in git stash. Use 'git stash pop' to recover them.")
		return
	}
	ui.Success(fmt.Sprintf("Switched to branch: %s", originalBranch))

	// 5. 删除本地新分支
	ui.Info("Deleting local branch...")
	if err := git.DeleteBranch(newBranch); err != nil {
		ui.Warning(fmt.Sprintf("Failed to delete local branch: %v", err))
	} else {
		ui.Success("Local branch deleted")
	}

	// 6. 恢复暂存的改动
	ui.Info("Restoring your changes...")
	stashPopCmd := exec.Command("git", "stash", "pop")
	if err := stashPopCmd.Run(); err != nil {
		ui.Warning(fmt.Sprintf("Failed to restore changes: %v", err))
		ui.Warning("Your changes are in git stash. Use 'git stash pop' to recover them.")
	} else {
		ui.Success("Changes restored to staging area")
	}

	ui.Success("Rollback completed! Your changes are preserved.")
}
