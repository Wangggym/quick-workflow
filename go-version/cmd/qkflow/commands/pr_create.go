package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Wangggym/quick-workflow/internal/ai"
	"github.com/Wangggym/quick-workflow/internal/editor"
	"github.com/Wangggym/quick-workflow/internal/git"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/watcher"
	"github.com/Wangggym/quick-workflow/internal/config"
	"github.com/Wangggym/quick-workflow/internal/ui"
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
		log.Error("Not a git repository")
		return
	}

	// 检查是否有未提交的更改
	hasChanges, err := git.HasUncommittedChanges()
	if err != nil {
		log.Error("Failed to check git status: %v", err)
		return
	}
	if !hasChanges {
		log.Error("No changes to commit. Please stage your changes first with 'git add'")
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
				log.Warning("Operation cancelled by user")
				os.Exit(0)
			}
			log.Error("Failed to get input: %v", err)
			return
		}
	}

	// 如果有 Jira ticket，获取信息
	var jiraIssue *jira.Issue
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		jiraClient, err := jira.NewClient()
		if err != nil {
			log.Warning("Failed to create Jira client: %v", err)
		} else {
			jiraIssue, err = jiraClient.GetIssue(jiraTicket)
			if err != nil {
				log.Warning("Failed to get Jira issue: %v", err)
			} else {
				log.Info("Found Jira issue: %s", jiraIssue.Summary)
			}
		}
	}

	// 显示 Jira 信息
	if jiraIssue != nil {
		log.Info("Jira issue: %s", jiraIssue.Summary)
	}

	// 选择变更类型
	prTypes := ui.PRTypeOptions()
	selectedTypes, err := ui.PromptMultiSelect("Select type(s) of changes:", prTypes)
	if err != nil {
		if err.Error() == "interrupt" {
			log.Warning("Operation cancelled by user")
			os.Exit(0)
		}
		log.Warning("No types selected, continuing...")
		selectedTypes = []string{}
	}

	// 询问是否添加说明/截图 (空格选择，Enter 跳过)
	var editorResult *editor.EditorResult
	addDescription, err := ui.PromptOptional("Add detailed description with images/videos?")
	if err == nil && addDescription {
		log.Info("Opening web editor...")
		editorResult, err = editor.StartEditor()
		if err != nil {
			log.Warning("Failed to start editor: %v", err)
			editorResult = nil
		} else if editorResult.Content == "" && len(editorResult.Files) == 0 {
			log.Info("No content added, skipping...")
			editorResult = nil
		} else {
			log.Success("Content saved! (%d characters, %d files)", len(editorResult.Content), len(editorResult.Files))
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
			log.Info("Generating PR title with AI...")
			title, err = aiClient.GeneratePRTitle(jiraIssue.Summary, prType, "")
			if err != nil {
				log.Warning("AI generation failed: %v", err)
				// 回退到简单格式
				title = generateSimpleTitle(jiraIssue.Summary, prType, "")
			} else {
				log.Success("Generated title: %s", title)
			}
		} else {
			// 没有 AI 或没有类型，使用简单生成
			title = generateSimpleTitle(jiraIssue.Summary, prType, "")
			log.Success("Generated title: %s", title)
		}
	} else {
		// 没有 Jira，手动输入
		title, err = ui.PromptInput("Enter PR title:", true)
		if err != nil {
			if err.Error() == "interrupt" {
				log.Warning("Operation cancelled by user")
				os.Exit(0)
			}
			log.Error("Failed to get title: %v", err)
			return
		}
	}

	// 构建 PR body
	prBody := buildPRBody(selectedTypes, jiraTicket)

	// 创建分支名
	branchName := buildBranchName(jiraTicket, title)
	cfg := config.Get()
	if cfg.BranchPrefix != "" {
		branchName = cfg.BranchPrefix + "/" + branchName
	}

	log.Info("Creating branch: %s", branchName)

	// 创建分支
	if err := git.CreateBranch(branchName); err != nil {
		log.Error("Failed to create branch: %v", err)
		return
	}

	// Stage 所有更改
	log.Info("Staging changes...")
	if err := git.AddAll(); err != nil {
		log.Error("Failed to stage changes: %v", err)
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

	log.Info("Committing changes...")
	if err := git.Commit(commitMessage); err != nil {
		log.Error("Failed to commit: %v", err)
		return
	}

	// 推送分支
	log.Info("Pushing branch to remote...")
	if err := git.Push(branchName); err != nil {
		log.Error("Failed to push: %v", err)
		return
	}

	// 获取仓库信息
	remoteURL, err := git.GetRemoteURL()
	if err != nil {
		log.Error("Failed to get remote URL: %v", err)
		return
	}

	owner, repo, err := github.ParseRepositoryFromURL(remoteURL)
	if err != nil {
		log.Error("Failed to parse repository: %v", err)
		return
	}

	// 获取默认分支
	defaultBranch, err := git.GetDefaultBranch()
	if err != nil {
		log.Warning("Failed to detect default branch, using 'main': %v", err)
		defaultBranch = "main"
	}
	log.Info("Using base branch: %s", defaultBranch)

	// 创建 PR
	log.Info("Creating pull request...")
	ghClient, err := github.NewClient()
	if err != nil {
		log.Error("Failed to create GitHub client: %v", err)
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
		log.Error("Failed to create PR: %v", err)
		return
	}

	log.Success("Pull request created: %s", pr.HTMLURL)

	// 处理编辑器内容（上传文件并添加评论）
	if editorResult != nil && (editorResult.Content != "" || len(editorResult.Files) > 0) {
		log.Info("Processing description and files...")

		// 创建 Jira 客户端（如果需要）
		var jiraClient *jira.Client
		if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
			jiraClient, err = jira.NewClient()
			if err != nil {
				log.Warning("Failed to create Jira client for file upload: %v", err)
				jiraClient = nil
			}
		}

		// 上传文件
		var uploadResults []editor.UploadResult
		if len(editorResult.Files) > 0 {
			log.Info("Uploading %d file(s)...", len(editorResult.Files))
			uploadResults, err = editor.UploadFiles(
				editorResult.Files,
				ghClient,
				jiraClient,
				pr.Number,
				owner,
				repo,
				jiraTicket,
			)
			if err != nil {
				log.Warning("Failed to upload files: %v", err)
			} else {
				log.Success("Uploaded %d file(s)", len(uploadResults))
			}
		}

		// 替换 markdown 中的本地路径为在线 URL
		content := editorResult.Content
		if len(uploadResults) > 0 {
			content = editor.ReplaceLocalPathsWithURLs(content, uploadResults)
		}

		// 添加评论到 GitHub PR
		if content != "" {
			log.Info("Adding description to GitHub PR...")
			if err := ghClient.AddPRComment(owner, repo, pr.Number, content); err != nil {
				log.Warning("Failed to add comment to GitHub: %v", err)
			} else {
				log.Success("Description added to GitHub PR")
			}
		}

		// 添加评论到 Jira
		if jiraClient != nil && jiraTicket != "" && content != "" {
			log.Info("Adding description to Jira...")
			jiraComment := fmt.Sprintf("*PR Description:*\n\n%s\n\n[View PR|%s]", content, pr.HTMLURL)
			if err := jiraClient.AddComment(jiraTicket, jiraComment); err != nil {
				log.Warning("Failed to add comment to Jira: %v", err)
			} else {
				log.Success("Description added to Jira")
			}
		}

		// 清理临时文件
		if len(editorResult.Files) > 0 {
			for _, file := range editorResult.Files {
				os.Remove(file)
			}
		}
	}

	// 更新 Jira
	if jiraTicket != "" && jira.ValidateIssueKey(jiraTicket) {
		jiraClient, err := jira.NewClient()
		if err != nil {
			log.Warning("Failed to create Jira client: %v", err)
		} else {
			// 分配给当前用户
			log.Info("Assigning Jira ticket to you...")
			if err := jiraClient.AssignToMe(jiraTicket); err != nil {
				log.Warning("Failed to assign ticket: %v", err)
			} else {
				log.Success("Assigned Jira ticket to you")
			}

			// 添加 PR 链接
			log.Info("Adding PR link to Jira...")
			if err := jiraClient.AddPRLink(jiraTicket, pr.HTMLURL); err != nil {
				log.Warning("Failed to add PR link to Jira: %v", err)
			} else {
				log.Success("Added PR link to Jira")
			}

			// 更新状态
				projectKey := jira.ExtractProjectKey(jiraTicket)

			// 检查状态缓存
			statusCache, err := jira.NewStatusCache()
			if err != nil {
				log.Warning("Failed to create status cache: %v", err)
			} else {
				mapping, err := statusCache.GetProjectStatus(projectKey)
				if err != nil {
					log.Warning("Failed to get cached status: %v", err)
				} else if mapping == nil {
					// 第一次使用，配置状态映射
					log.Info("First time using project %s, please configure status mappings", projectKey)
					mapping, err = setupProjectStatusMapping(jiraClient, projectKey)
					if err != nil {
						log.Warning("Failed to setup status mapping: %v", err)
					} else if mapping != nil {
						// 保存配置
						if err := statusCache.SaveProjectStatus(mapping); err != nil {
							log.Warning("Failed to save status mapping: %v", err)
						} else {
							log.Success("Status mapping saved!")
						}
					}
				}

				// 使用缓存的状态更新
				if mapping != nil && mapping.PRCreatedStatus != "" {
					log.Info("Updating Jira status to: %s", mapping.PRCreatedStatus)
					if err := jiraClient.UpdateStatus(jiraTicket, mapping.PRCreatedStatus); err != nil {
						log.Warning("Failed to update status: %v", err)
					} else {
						log.Success("Updated Jira status to: %s", mapping.PRCreatedStatus)
					}
				}
			}
		}
	}

	// 添加到 watching list
	watchingList, err := watcher.NewWatchingList()
	if err != nil {
		log.Warning("Failed to load watching list: %v", err)
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
			log.Warning("Failed to add PR to watching list: %v", err)
		} else {
			log.Info("✅ Added PR to watching list for auto Jira updates")
		}
	}

	// 复制 URL 到剪贴板
	copyToClipboard(pr.HTMLURL)

	// 打开浏览器
	openBrowser(pr.HTMLURL)

	log.Info("")
	log.Success("All done! 🎉")
}

func buildBranchName(jiraTicket, title string) string {
	sanitized := git.SanitizeBranchName(title)
	if jiraTicket != "" {
		return fmt.Sprintf("%s--%s", jiraTicket, sanitized)
	}
	return sanitized
}

func buildPRBody(types []string, jiraTicket string) string {
	var body strings.Builder

	body.WriteString("# PR Ready\n\n")

	if len(types) > 0 {
		body.WriteString("## Types of changes\n\n")
		for _, t := range types {
			body.WriteString(fmt.Sprintf("- [x] %s\n", t))
		}
		body.WriteString("\n")
	}

	if jiraTicket != "" {
		cfg := config.Get()
		jiraURL := fmt.Sprintf("%s/browse/%s", cfg.JiraServiceAddress, jiraTicket)
		body.WriteString(fmt.Sprintf("#### Jira Link:\n\n%s\n", jiraURL))
	}

	return body.String()
}

func setupProjectStatusMapping(client *jira.Client, projectKey string) (*jira.StatusMapping, error) {
	statuses, err := client.GetProjectStatuses(projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get project statuses: %w", err)
	}

	log.Info("Select status when PR is created/in progress:")
	createdStatus, err := ui.PromptSelect("Status for PR created:", statuses)
	if err != nil {
		if err.Error() == "interrupt" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select created status: %w", err)
	}

	log.Info("Select status when PR is merged/done:")
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

func generateSimpleTitle(jiraSummary, prType, description string) string {
	// 如果有简短描述，使用描述
	if description != "" {
		if prType != "" {
			return fmt.Sprintf("%s: %s", prType, description)
		}
		return description
	}

	// 否则使用 Jira 标题的前 50 个字符
	summary := jiraSummary
	if len(summary) > 50 {
		summary = summary[:50] + "..."
	}

	if prType != "" {
		return fmt.Sprintf("%s: %s", prType, summary)
	}
	return summary
}

func copyToClipboard(text string) {
	// macOS only for now
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		log.Warning("Failed to copy to clipboard")
	} else {
		log.Success("Successfully copied %s to clipboard", text)
	}
}

func openBrowser(url string) {
	// macOS
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		log.Warning("Failed to open browser: %v", err)
	}
}

