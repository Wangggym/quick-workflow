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
			ui.Info("Detected GitHub PR URL, parsing...")
			owner, repo, prNumber, err = github.ParsePRFromURL(arg)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to parse PR URL: %v", err))
				return
			}
			ui.Success(fmt.Sprintf("Parsed: %s/%s PR #%d", owner, repo, prNumber))
		} else {
			// 尝试作为 PR 号解析
			prNumber, err = strconv.Atoi(arg)
			if err != nil {
				ui.Error(fmt.Sprintf("Invalid PR number or URL: %s", arg))
				ui.Info("Expected: PR number (e.g., '123') or GitHub URL (e.g., 'https://github.com/owner/repo/pull/123')")
				return
			}
			
			// PR 号格式，需要从本地仓库获取 owner/repo
			if !git.IsGitRepository() {
				ui.Error("Not a git repository. When using PR number, you must be in a git repository.")
				ui.Info("Alternatively, use the full GitHub PR URL: https://github.com/owner/repo/pull/NUMBER")
				return
			}
			
			remoteURL, err := git.GetRemoteURL()
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to get remote URL: %v", err))
				return
			}
			
			owner, repo, err = github.ParseRepositoryFromURL(remoteURL)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to parse repository: %v", err))
				return
			}
		}
	} else {
		// 没有提供参数，使用原有的自动检测逻辑
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

		owner, repo, err = github.ParseRepositoryFromURL(remoteURL)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to parse repository: %v", err))
			return
		}
	}

	// 创建 GitHub 客户端
	ghClient, err := github.NewClient()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create GitHub client: %v", err))
		return
	}

	// 如果还没有 PR 号，尝试自动检测或让用户选择
	if prNumber == 0 {
		// 尝试获取当前分支的 PR
		currentBranch, err := git.GetCurrentBranch()
		if err == nil && currentBranch != "" {
			ui.Info(fmt.Sprintf("Checking for PR from current branch: %s", currentBranch))

			// 先尝试 open 状态的 PR
			prs, err := ghClient.ListPullRequests(owner, repo, "open", "")
			if err == nil {
				for _, pr := range prs {
					if pr.Head == currentBranch {
						prNumber = pr.Number
						ui.Success(fmt.Sprintf("Found PR #%d: %s", pr.Number, pr.Title))
						break
					}
				}
			}

			// 如果还是没找到，提示用户
			if prNumber == 0 {
				ui.Warning(fmt.Sprintf("No PR found for branch: %s", currentBranch))
				fmt.Println()

				// 询问用户是否从列表选择
				manually, err := ui.PromptConfirm("Do you want to select a PR from the list?", true)
				if err != nil || !manually {
					ui.Info("Approve cancelled")
					return
				}
			}
		}

		// 如果用户选择手动选择或没有当前分支
		if prNumber == 0 {
			prs, err := ghClient.ListPullRequests(owner, repo, "open", "")
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

			selected, err := ui.PromptSelect("Select a PR to approve:", prOptions)
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
	if pr.State != "open" {
		ui.Error(fmt.Sprintf("PR is not open (state: %s)", pr.State))
		return
	}

	// 如果没有提供评论，使用默认值 "👍"
	comment := approveComment
	if comment == "" {
		comment = "👍"
		ui.Info("Using default comment: 👍 (use -c flag to customize)")
	}

	// 批准 PR
	ui.Info(fmt.Sprintf("Approving PR #%d...", prNumber))
	approvalSucceeded := true
	if err := ghClient.ApprovePullRequest(owner, repo, prNumber, comment); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "422") {
			approvalSucceeded = false
			ui.Warning("Cannot approve this PR (you may be the author or already approved)")
			fmt.Println()
			
			// 如果带了 -m 参数，直接跳过批准继续合并
			if approveAndMerge {
				ui.Info("💡 Skipping approval, proceeding directly to merge...")
			} else {
				// 没有 -m 参数，提示错误并退出
				ui.Error("Approval failed. If you want to merge directly, use the -m flag:")
				ui.Info(fmt.Sprintf("  qkflow pr approve %d -m", prNumber))
				return
			}
		} else {
			ui.Error(fmt.Sprintf("Failed to approve PR: %v", err))
			return
		}
	}

	if approvalSucceeded {
		if comment != "" {
			ui.Success(fmt.Sprintf("✅ PR approved with comment: %s", comment))
		} else {
			ui.Success("✅ PR approved!")
		}
	}

	// 如果需要自动合并
	if approveAndMerge {
		// 检查 PR 是否可以合并
		ui.Info("Checking if PR is mergeable...")
		mergeable, err := ghClient.IsPRMergeable(owner, repo, prNumber)
		if err != nil || !mergeable {
			ui.Warning(fmt.Sprintf("Cannot merge PR: %v", err))
			ui.Info("You may need to wait for CI checks or resolve conflicts")
			return
		}

		// 执行合并
		ui.Info(fmt.Sprintf("Merging PR #%d...", prNumber))
		if err := ghClient.MergePullRequest(owner, repo, prNumber, pr.Title); err != nil {
			ui.Error(fmt.Sprintf("Failed to merge PR: %v", err))
			return
		}

		ui.Success("🎉 PR merged successfully!")

		// 删除远程分支
		ui.Info(fmt.Sprintf("Deleting remote branch %s...", pr.Head))
		if err := git.DeleteRemoteBranch(pr.Head); err != nil {
			ui.Warning(fmt.Sprintf("Failed to delete remote branch: %v (may already be deleted)", err))
		} else {
			ui.Success("Remote branch deleted")
		}

		// 如果在同一个分支，切换到主分支
		currentBranch, err := git.GetCurrentBranch()
		if err == nil && currentBranch == pr.Head {
			defaultBranch, err := git.GetDefaultBranch()
			if err != nil {
				defaultBranch = "master"
			}

			ui.Info(fmt.Sprintf("Switching to %s branch...", defaultBranch))
			if err := git.CheckoutBranch(defaultBranch); err != nil {
				ui.Warning(fmt.Sprintf("Could not switch to %s, you may need to do this manually", defaultBranch))
			} else {
				ui.Success("Switched to default branch")

				// 删除本地分支
				ui.Info(fmt.Sprintf("Deleting local branch %s...", pr.Head))
				if err := git.DeleteBranch(pr.Head); err != nil {
					ui.Warning(fmt.Sprintf("Failed to delete local branch: %v", err))
				} else {
					ui.Success("Local branch deleted")
				}
			}
		}

		fmt.Println()
		ui.Success("All done! 🎉")
	} else {
		fmt.Println()
		ui.Info("PR approved. Use 'qkg pr merge' to merge it later, or run with --merge flag to auto-merge.")
	}
}

