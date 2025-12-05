package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Wangggym/quick-workflow/internal/git"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/spf13/cobra"
)

var (
	approveAndMerge bool
	approveComment  string
)

var prApproveCmd = &cobra.Command{
	Use:   "approve [pr-number|pr-url]",
	Short: "Approve a PR and optionally merge it",
	Long: `Approve a pull request and optionally merge it automatically:
  - Approve the PR on GitHub
  - Add a comment (default: 👍, customize with -c flag)
  - Optionally auto-merge after approval

Arguments:
  [pr-number|pr-url]  PR number (e.g., 123) or full GitHub PR URL
                      (e.g., https://github.com/owner/repo/pull/123)
                      Omit to auto-detect from current branch

Examples:
  qkflow pr approve 123                  # Approves with 👍
  qkflow pr approve 123 -c "LGTM!"      # Custom comment
  qkflow pr approve https://github.com/brain/planning-api/pull/2001
  qkflow pr approve 123 -m               # Approve with 👍 and merge`,
	Args: cobra.MaximumNArgs(1),
	Run:  runPRApprove,
}

func init() {
	prApproveCmd.Flags().BoolVarP(&approveAndMerge, "merge", "m", false, "Automatically merge the PR after approval")
	prApproveCmd.Flags().StringVarP(&approveComment, "comment", "c", "", "Add a comment with the approval (default: 👍)")
}

func runPRApprove(cmd *cobra.Command, args []string) {
	var owner, repo string
	var prNumber int
	var err error

	// 如果提供了参数，检查是否是 URL 格式
	if len(args) > 0 {
		arg := args[0]

		// 检查是否是 GitHub PR URL
		if github.IsPRURL(arg) {
			log.Info("Detected GitHub PR URL, parsing...")
			owner, repo, prNumber, err = github.ParsePRFromURL(arg)
			if err != nil {
				log.Error("Failed to parse PR URL: %v", err)
				return
			}
			log.Success("Parsed: %s/%s PR #%d", owner, repo, prNumber)
		} else {
			// 尝试作为 PR 号解析
			prNumber, err = strconv.Atoi(arg)
			if err != nil {
				log.Error("Invalid PR number or URL: %s", arg)
				log.Info("Expected: PR number (e.g., '123') or GitHub URL (e.g., 'https://github.com/owner/repo/pull/123')")
				return
			}

			// PR 号格式，需要从本地仓库获取 owner/repo
			if !git.IsGitRepository() {
				log.Error("Not a git repository. When using PR number, you must be in a git repository.")
				log.Info("Alternatively, use the full GitHub PR URL: https://github.com/owner/repo/pull/NUMBER")
				return
			}

			remoteURL, err := git.GetRemoteURL()
			if err != nil {
				log.Error("Failed to get remote URL: %v", err)
				return
			}

			owner, repo, err = github.ParseRepositoryFromURL(remoteURL)
			if err != nil {
				log.Error("Failed to parse repository: %v", err)
				return
			}
		}
	} else {
		// 没有提供参数，使用原有的自动检测逻辑
		// 检查是否在 Git 仓库中
		if !git.IsGitRepository() {
			log.Error("Not a git repository")
			return
		}

		// 获取仓库信息
		remoteURL, err := git.GetRemoteURL()
		if err != nil {
			log.Error("Failed to get remote URL: %v", err)
			return
		}

		owner, repo, err = github.ParseRepositoryFromURL(remoteURL)
		if err != nil {
			log.Error("Failed to parse repository: %v", err)
			return
		}
	}

	// 创建 GitHub 客户端
	ghClient, err := github.NewClient()
	if err != nil {
		log.Error("Failed to create GitHub client: %v", err)
		return
	}

	// 如果还没有 PR 号，尝试自动检测或让用户选择
	if prNumber == 0 {
		// 尝试获取当前分支的 PR
		currentBranch, err := git.GetCurrentBranch()
		if err == nil && currentBranch != "" {
			log.Info("Checking for PR from current branch: %s", currentBranch)

			// 先尝试 open 状态的 PR
			prs, err := ghClient.ListPullRequests(owner, repo, "open", "")
			if err == nil {
				for _, pr := range prs {
					if pr.Head == currentBranch {
						prNumber = pr.Number
						log.Success("Found PR #%d: %s", pr.Number, pr.Title)
						break
					}
				}
			}

			// 如果还是没找到，提示用户
			if prNumber == 0 {
				log.Warning("No PR found for branch: %s", currentBranch)
				log.Info("")

				// 询问用户是否从列表选择
				manually, err := ui.PromptConfirm("Do you want to select a PR from the list?", true)
				if err != nil || !manually {
					log.Info("Approve cancelled")
					return
				}
			}
		}

		// 如果用户选择手动选择或没有当前分支
		if prNumber == 0 {
			prs, err := ghClient.ListPullRequests(owner, repo, "open", "")
			if err != nil {
				log.Error("Failed to list PRs: %v", err)
				return
			}

			if len(prs) == 0 {
				log.Error("No open pull requests found")
				return
			}

			// 构建选择列表
			prOptions := make([]string, len(prs))
			for i, pr := range prs {
				prOptions[i] = fmt.Sprintf("#%d - %s", pr.Number, pr.Title)
			}

			selected, err := ui.PromptSelect("Select a PR to approve:", prOptions)
			if err != nil {
				if err.Error() == "interrupt" {
					log.Warning("Operation cancelled by user")
					os.Exit(0)
				}
				log.Error("Failed to select PR: %v", err)
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
				log.Error("Failed to find selected PR")
				return
			}
		}
	}

	// 获取 PR 信息
	log.Info("Fetching PR #%d...", prNumber)
	pr, err := ghClient.GetPullRequest(owner, repo, prNumber)
	if err != nil {
		log.Error("Failed to get PR: %v", err)
		return
	}

	log.Info("PR: %s", pr.Title)
	log.Info("Branch: %s -> %s", pr.Head, pr.Base)
	log.Info("State: %s", pr.State)

	// 检查 PR 状态
	if pr.State != "open" {
		log.Error("PR is not open (state: %s)", pr.State)
		return
	}

	// 如果没有提供评论，使用默认值 "👍"
	comment := approveComment
	if comment == "" {
		comment = "👍"
		log.Info("Using default comment: 👍 (use -c flag to customize)")
	}

	// 批准 PR
	log.Info("Approving PR #%d...", prNumber)
	approvalSucceeded := true
	if err := ghClient.ApprovePullRequest(owner, repo, prNumber, comment); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "422") {
			approvalSucceeded = false
			log.Warning("Cannot approve this PR (you may be the author or already approved)")
			log.Info("")

			// 如果带了 -m 参数，直接跳过批准继续合并
			if approveAndMerge {
				log.Info("💡 Skipping approval, proceeding directly to merge...")
			} else {
				// 没有 -m 参数，提示错误并退出
				log.Error("Approval failed. If you want to merge directly, use the -m flag:")
				log.Info("  qkflow pr approve %d -m", prNumber)
				return
			}
		} else {
			log.Error("Failed to approve PR: %v", err)
			return
		}
	}

	if approvalSucceeded {
		if comment != "" {
			log.Success("✅ PR approved with comment: %s", comment)
		} else {
			log.Success("✅ PR approved!")
		}
	}

	// 如果需要自动合并
	if approveAndMerge {
		// 检查 PR 是否可以合并
		log.Info("Checking if PR is mergeable...")
		mergeable, err := ghClient.IsPRMergeable(owner, repo, prNumber)
		if err != nil || !mergeable {
			log.Warning("Cannot merge PR: %v", err)
			log.Info("You may need to wait for CI checks or resolve conflicts")
			return
		}

		// 执行合并
		log.Info("Merging PR #%d...", prNumber)
		if err := ghClient.MergePullRequest(owner, repo, prNumber, pr.Title); err != nil {
			log.Error("Failed to merge PR: %v", err)
			return
		}

		log.Success("🎉 PR merged successfully!")

		// 删除远程分支
		log.Info("Deleting remote branch %s...", pr.Head)
		if err := git.DeleteRemoteBranch(pr.Head); err != nil {
			log.Warning("Failed to delete remote branch: %v (may already be deleted)", err)
		} else {
			log.Success("Remote branch deleted")
		}

		// 如果在同一个分支，切换到主分支
		currentBranch, err := git.GetCurrentBranch()
		if err == nil && currentBranch == pr.Head {
			defaultBranch, err := git.GetDefaultBranch()
			if err != nil {
				defaultBranch = "master"
			}

			log.Info("Switching to %s branch...", defaultBranch)
			if err := git.CheckoutBranch(defaultBranch); err != nil {
				log.Warning("Could not switch to %s, you may need to do this manually", defaultBranch)
			} else {
				log.Success("Switched to default branch")

				// 删除本地分支
				log.Info("Deleting local branch %s...", pr.Head)
				if err := git.DeleteBranch(pr.Head); err != nil {
					log.Warning("Failed to delete local branch: %v", err)
				} else {
					log.Success("Local branch deleted")
				}
			}
		}

		log.Info("")
		log.Success("All done! 🎉")
	} else {
		log.Info("")
		log.Info("PR approved. Use 'qkg pr merge' to merge it later, or run with --merge flag to auto-merge.")
	}
}
