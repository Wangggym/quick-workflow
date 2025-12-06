package commands

import (

	"github.com/Wangggym/quick-workflow/internal/git"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Quick update - commit and push with PR title as commit message",
	Long: `Automatically get the current PR title and use it as the commit message.
This command will:
1. Get the PR title for the current branch
2. Stage all changes (git add --all)
3. Commit with the PR title as the message
4. Push to origin

If no PR is found, it will use "update" as the default commit message.`,
	Run: runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) {
	// 检查是否是 git 仓库
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
		log.Warning("No changes to commit")
		return
	}

	// 获取当前分支
	branch, err := git.GetCurrentBranch()
	if err != nil {
		log.Error("Failed to get current branch: %v", err)
		return
	}

	log.Info("Current branch: %s", branch)

	// 获取 PR 标题作为 commit message
	commitMessage := "update" // 默认 commit message
	
	// 尝试从 GitHub 获取 PR 标题
	remoteURL, err := git.GetRemoteURL()
	if err == nil {
		owner, repo, err := github.ParseRepositoryFromURL(remoteURL)
		if err == nil {
			// 创建 GitHub 客户端
			ghClient, err := github.NewClient()
			if err == nil {
				// 尝试获取当前分支的 PR
				pr, err := ghClient.GetPRByBranch(owner, repo, branch)
				if err == nil && pr != nil {
					commitMessage = pr.Title
					log.Success("Got PR title: %s", commitMessage)
				} else {
					log.Warning("No open PR found for branch %s, using default message 'update'", branch)
				}
			} else {
				log.Warning("Failed to create GitHub client: %v, using default message", err)
			}
		}
	}

	// Stage all changes
	log.Info("Staging all changes...")
	if err := git.AddAll(); err != nil {
		log.Error("Failed to stage changes: %v", err)
		return
	}

	// Commit
	log.Info("Committing with message: '%s'", commitMessage)
	if err := git.Commit(commitMessage); err != nil {
		log.Error("Failed to commit: %v", err)
		return
	}

	// Push
	log.Info("Pushing to origin...")
	if err := git.Push(branch); err != nil {
		log.Error("Failed to push: %v", err)
		return
	}

	log.Success("✅ Successfully committed and pushed changes!")
}

