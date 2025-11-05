package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/Wangggym/quick-workflow/pkg/config"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates a new GitHub client
func NewClient() (*Client, error) {
	cfg := config.Get()
	if cfg == nil || cfg.GitHubToken == "" {
		return nil, fmt.Errorf("GitHub token not configured")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}, nil
}

// PullRequest represents a pull request
type PullRequest struct {
	Number  int
	Title   string
	Body    string
	HTMLURL string
	Head    string
	Base    string
	State   string
}

// CreatePullRequestInput contains the input for creating a PR
type CreatePullRequestInput struct {
	Owner string
	Repo  string
	Title string
	Body  string
	Head  string // branch name
	Base  string // target branch, usually "main" or "master"
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(input CreatePullRequestInput) (*PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(input.Title),
		Body:  github.String(input.Body),
		Head:  github.String(input.Head),
		Base:  github.String(input.Base),
	}

	pr, _, err := c.client.PullRequests.Create(c.ctx, input.Owner, input.Repo, newPR)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return &PullRequest{
		Number:  pr.GetNumber(),
		Title:   pr.GetTitle(),
		Body:    pr.GetBody(),
		HTMLURL: pr.GetHTMLURL(),
		Head:    pr.GetHead().GetRef(),
		Base:    pr.GetBase().GetRef(),
		State:   pr.GetState(),
	}, nil
}

// GetPullRequest gets a pull request by number
func (c *Client) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(c.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &PullRequest{
		Number:  pr.GetNumber(),
		Title:   pr.GetTitle(),
		Body:    pr.GetBody(),
		HTMLURL: pr.GetHTMLURL(),
		Head:    pr.GetHead().GetRef(),
		Base:    pr.GetBase().GetRef(),
		State:   pr.GetState(),
	}, nil
}

// MergePullRequest merges a pull request
func (c *Client) MergePullRequest(owner, repo string, number int, commitMessage string) error {
	options := &github.PullRequestOptions{
		CommitTitle: commitMessage,
	}

	_, _, err := c.client.PullRequests.Merge(c.ctx, owner, repo, number, commitMessage, options)
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}

	return nil
}

// GetCurrentRepository gets the owner and repo from git remote
func GetCurrentRepository() (owner, repo string, err error) {
	// 这里可以通过执行 git remote get-url origin 来获取
	// 简化实现：从环境变量或配置读取
	// 实际应该解析 git remote 输出
	
	// TODO: 实现通过 git 命令获取
	return "", "", fmt.Errorf("not implemented: should parse git remote")
}

// ParseRepositoryFromURL parses owner and repo from GitHub URL
func ParseRepositoryFromURL(url string) (owner, repo string, err error) {
	// 支持多种格式：
	// https://github.com/owner/repo
	// git@github.com:owner/repo.git
	// owner/repo

	url = strings.TrimSpace(url)

	// 移除协议部分
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "git@")

	// 移除 github.com
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimPrefix(url, "github.com:")

	// 移除 .git 后缀
	url = strings.TrimSuffix(url, ".git")

	// 分割 owner/repo
	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", url)
	}

	return parts[0], parts[1], nil
}

