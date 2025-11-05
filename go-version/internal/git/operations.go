package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CheckStatus checks if there are uncommitted changes
func CheckStatus() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	// 如果有输出，说明有未提交的更改
	return len(output) > 0, nil
}

// GetCurrentBranch gets the current branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateBranch creates and checks out a new branch
func CreateBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch %s: %w\n%s", branchName, err, stderr.String())
	}

	return nil
}

// AddAll stages all changes
func AddAll() error {
	cmd := exec.Command("git", "add", "--all")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w\n%s", err, stderr.String())
	}

	return nil
}

// Commit creates a commit with the given message
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message, "--no-verify")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w\n%s", err, stderr.String())
	}

	return nil
}

// Push pushes the current branch to origin
func Push(branchName string) error {
	cmd := exec.Command("git", "push", "-u", "origin", branchName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch %s: %w\n%s", branchName, err, stderr.String())
	}

	return nil
}

// DeleteBranch deletes a local branch
func DeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-d", branchName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete branch %s: %w\n%s", branchName, err, stderr.String())
	}

	return nil
}

// DeleteRemoteBranch deletes a remote branch
func DeleteRemoteBranch(branchName string) error {
	cmd := exec.Command("git", "push", "origin", "--delete", branchName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete remote branch %s: %w\n%s", branchName, err, stderr.String())
	}

	return nil
}

// GetRemoteURL gets the remote URL
func GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// SanitizeBranchName converts a string to a valid branch name
func SanitizeBranchName(name string) string {
	// 替换非法字符为连字符
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, name)

	// 移除多余的连字符
	name = strings.ReplaceAll(name, "--", "-")
	name = strings.Trim(name, "-")

	return name
}

// HasUncommittedChanges checks if there are uncommitted changes
func HasUncommittedChanges() (bool, error) {
	return CheckStatus()
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetDefaultBranch gets the default branch name (usually main or master)
func GetDefaultBranch() (string, error) {
	// 尝试从 origin/HEAD 获取默认分支
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		// 输出格式: refs/remotes/origin/main
		branch := strings.TrimSpace(string(output))
		branch = strings.TrimPrefix(branch, "refs/remotes/origin/")
		if branch != "" {
			return branch, nil
		}
	}

	// 如果上面失败，尝试列出远程分支并检测 main 或 master
	cmd = exec.Command("git", "remote", "show", "origin")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "HEAD branch:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					branch := strings.TrimSpace(parts[1])
					if branch != "" {
						return branch, nil
					}
				}
			}
		}
	}

	// 最后的回退方案：检查 main 或 master 是否存在
	for _, branch := range []string{"main", "master"} {
		cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/"+branch)
		if cmd.Run() == nil {
			return branch, nil
		}
	}

	// 默认返回 main
	return "main", nil
}

