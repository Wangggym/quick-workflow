package jira

import (
	"fmt"
	"strings"

	"github.com/Wangggym/quick-workflow/pkg/config"
	jira "github.com/andygrunwald/go-jira"
)

// Client wraps the Jira API client
type Client struct {
	client *jira.Client
}

// NewClient creates a new Jira client
func NewClient() (*Client, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	if cfg.JiraAPIToken == "" || cfg.JiraServiceAddress == "" || cfg.Email == "" {
		return nil, fmt.Errorf("Jira credentials not configured")
	}

	tp := jira.BasicAuthTransport{
		Username: cfg.Email,
		Password: cfg.JiraAPIToken,
	}

	client, err := jira.NewClient(tp.Client(), cfg.JiraServiceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return &Client{client: client}, nil
}

// Issue represents a Jira issue
type Issue struct {
	Key         string
	Summary     string
	Description string
	Status      string
	Type        string
}

// GetIssue gets a Jira issue by key
func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	issue, _, err := c.client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %s: %w", issueKey, err)
	}

	return &Issue{
		Key:         issue.Key,
		Summary:     issue.Fields.Summary,
		Description: issue.Fields.Description,
		Status:      issue.Fields.Status.Name,
		Type:        issue.Fields.Type.Name,
	}, nil
}

// UpdateStatus updates the status of a Jira issue
func (c *Client) UpdateStatus(issueKey, newStatus string) error {
	// 获取可用的转换
	transitions, _, err := c.client.Issue.GetTransitions(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get transitions for %s: %w", issueKey, err)
	}

	// 查找匹配的转换
	var targetTransition *jira.Transition
	for _, t := range transitions {
		if strings.EqualFold(t.Name, newStatus) || strings.EqualFold(t.To.Name, newStatus) {
			targetTransition = &t
			break
		}
	}

	if targetTransition == nil {
		return fmt.Errorf("transition to status '%s' not found for issue %s", newStatus, issueKey)
	}

	// 执行转换
	_, err = c.client.Issue.DoTransition(issueKey, targetTransition.ID)
	if err != nil {
		return fmt.Errorf("failed to transition issue %s to %s: %w", issueKey, newStatus, err)
	}

	return nil
}

// AddComment adds a comment to a Jira issue
func (c *Client) AddComment(issueKey, comment string) error {
	jiraComment := &jira.Comment{
		Body: comment,
	}

	_, _, err := c.client.Issue.AddComment(issueKey, jiraComment)
	if err != nil {
		return fmt.Errorf("failed to add comment to %s: %w", issueKey, err)
	}

	return nil
}

// GetProjectStatuses gets all available statuses for a project
func (c *Client) GetProjectStatuses(projectKey string) ([]string, error) {
	// 获取项目信息
	project, _, err := c.client.Project.Get(projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", projectKey, err)
	}

	// 获取状态列表
	statuses, _, err := c.client.Status.GetAllStatuses()
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}

	var statusNames []string
	for _, status := range statuses {
		statusNames = append(statusNames, status.Name)
	}

	// 如果项目有特定的工作流，这里应该过滤
	// 简化实现：返回所有状态
	_ = project // 避免未使用变量警告
	
	return statusNames, nil
}

// AddPRLink adds a PR link to a Jira issue
func (c *Client) AddPRLink(issueKey, prURL string) error {
	comment := fmt.Sprintf("Pull Request created: %s", prURL)
	return c.AddComment(issueKey, comment)
}

// ExtractProjectKey extracts the project key from an issue key
// For example: "PROJ-123" -> "PROJ"
func ExtractProjectKey(issueKey string) string {
	parts := strings.Split(issueKey, "-")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// ValidateIssueKey checks if an issue key is valid
func ValidateIssueKey(issueKey string) bool {
	if issueKey == "" {
		return false
	}
	parts := strings.Split(issueKey, "-")
	return len(parts) >= 2
}

