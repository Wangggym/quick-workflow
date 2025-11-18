package jira

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Formatter handles formatting Jira issues for display
type Formatter struct {
	client *Client
}

// NewFormatter creates a new formatter
func NewFormatter(client *Client) *Formatter {
	return &Formatter{client: client}
}

// FormatIssueSimple formats an issue for simple display (show command)
func (f *Formatter) FormatIssueSimple(issue *Issue) string {
	var sb strings.Builder

	// Header
	sb.WriteString("╔═══════════════════════════════════════════════════════╗\n")
	sb.WriteString(fmt.Sprintf("║ 🎫 %s: %s\n", issue.Key, truncate(issue.Summary, 45)))
	sb.WriteString("╚═══════════════════════════════════════════════════════╝\n\n")

	// Metadata
	sb.WriteString(fmt.Sprintf("📋 Type:        %s\n", issue.Type))
	sb.WriteString(fmt.Sprintf("📊 Status:      %s\n", issue.Status))
	if issue.Priority != "" {
		sb.WriteString(fmt.Sprintf("🏷️  Priority:    %s\n", issue.Priority))
	}
	if issue.Assignee != "" {
		sb.WriteString(fmt.Sprintf("👤 Assignee:    %s\n", issue.Assignee))
	}
	if issue.Reporter != "" {
		sb.WriteString(fmt.Sprintf("📝 Reporter:    %s\n", issue.Reporter))
	}

	// URL
	if f.client != nil {
		url := f.client.GetJiraURL(issue.Key)
		if url != "" {
			sb.WriteString(fmt.Sprintf("\n🔗 View in Jira: %s\n", url))
		}
	}

	return sb.String()
}

// FormatIssueFull formats an issue with full details (show --full command)
func (f *Formatter) FormatIssueFull(issue *Issue) string {
	var sb strings.Builder

	// Header
	sb.WriteString("╔═══════════════════════════════════════════════════════╗\n")
	sb.WriteString(fmt.Sprintf("║ 🎫 %s: %s\n", issue.Key, truncate(issue.Summary, 45)))
	sb.WriteString("╚═══════════════════════════════════════════════════════╝\n\n")

	// Metadata
	sb.WriteString(fmt.Sprintf("📋 Type:        %s\n", issue.Type))
	sb.WriteString(fmt.Sprintf("📊 Status:      %s\n", issue.Status))
	if issue.Priority != "" {
		sb.WriteString(fmt.Sprintf("🏷️  Priority:    %s\n", issue.Priority))
	}
	if issue.Assignee != "" {
		sb.WriteString(fmt.Sprintf("👤 Assignee:    %s\n", issue.Assignee))
	}
	if issue.Reporter != "" {
		sb.WriteString(fmt.Sprintf("📝 Reporter:    %s\n", issue.Reporter))
	}
	if issue.Created != "" {
		sb.WriteString(fmt.Sprintf("📅 Created:     %s\n", formatDate(issue.Created)))
	}
	if issue.Updated != "" {
		sb.WriteString(fmt.Sprintf("🔄 Updated:     %s\n", formatDate(issue.Updated)))
	}

	// Description
	sb.WriteString("\n📝 Description:\n")
	sb.WriteString("──────────────────────────────────────────────────────\n")
	if issue.Description != "" {
		sb.WriteString(issue.Description)
		sb.WriteString("\n")
	} else {
		sb.WriteString("(No description)\n")
	}
	sb.WriteString("──────────────────────────────────────────────────────\n")

	// Attachments
	if len(issue.Attachments) > 0 {
		sb.WriteString(fmt.Sprintf("\n📎 Attachments: (%d)\n", len(issue.Attachments)))
		for _, att := range issue.Attachments {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", att.Filename, formatSize(att.Size)))
		}
		sb.WriteString("\n⚠️  Images available but not shown. Use 'export --with-images' to download.\n")
	}

	// Comments
	if len(issue.Comments) > 0 {
		sb.WriteString(fmt.Sprintf("\n💬 Comments: (%d)\n", len(issue.Comments)))
		// Show up to 5 most recent comments
		displayCount := len(issue.Comments)
		if displayCount > 5 {
			displayCount = 5
		}
		for i := len(issue.Comments) - displayCount; i < len(issue.Comments); i++ {
			c := issue.Comments[i]
			sb.WriteString(fmt.Sprintf("\n[%s - %s]\n", c.Author, formatDate(c.Created)))
			sb.WriteString(fmt.Sprintf("%s\n", c.Body))
		}
		if len(issue.Comments) > displayCount {
			sb.WriteString(fmt.Sprintf("\n... and %d more comments\n", len(issue.Comments)-displayCount))
		}
	}

	// URL
	if f.client != nil {
		url := f.client.GetJiraURL(issue.Key)
		if url != "" {
			sb.WriteString(fmt.Sprintf("\n🔗 View in Jira: %s\n", url))
		}
	}

	return sb.String()
}

// FormatMarkdown formats an issue as Markdown for export
func (f *Formatter) FormatMarkdown(issue *Issue, exportPath string, withImages bool) string {
	var sb strings.Builder

	// Front matter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("issue_key: %s\n", issue.Key))
	sb.WriteString(fmt.Sprintf("title: %s\n", issue.Summary))
	sb.WriteString(fmt.Sprintf("type: %s\n", issue.Type))
	sb.WriteString(fmt.Sprintf("status: %s\n", issue.Status))
	if issue.Priority != "" {
		sb.WriteString(fmt.Sprintf("priority: %s\n", issue.Priority))
	}
	if issue.Assignee != "" {
		sb.WriteString(fmt.Sprintf("assignee: %s\n", issue.Assignee))
	}
	if issue.Reporter != "" {
		sb.WriteString(fmt.Sprintf("reporter: %s\n", issue.Reporter))
	}
	if issue.Created != "" {
		sb.WriteString(fmt.Sprintf("created: %s\n", issue.Created))
	}
	if issue.Updated != "" {
		sb.WriteString(fmt.Sprintf("updated: %s\n", issue.Updated))
	}
	if f.client != nil {
		url := f.client.GetJiraURL(issue.Key)
		if url != "" {
			sb.WriteString(fmt.Sprintf("jira_url: %s\n", url))
		}
	}
	sb.WriteString(fmt.Sprintf("exported_at: %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("export_includes_images: %t\n", withImages))
	sb.WriteString("---\n\n")

	// Title
	sb.WriteString(fmt.Sprintf("# %s: %s\n\n", issue.Key, issue.Summary))

	// Metadata table
	sb.WriteString("## 📊 Metadata\n\n")
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|-------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **Type** | %s |\n", issue.Type))
	sb.WriteString(fmt.Sprintf("| **Status** | %s |\n", issue.Status))
	if issue.Priority != "" {
		sb.WriteString(fmt.Sprintf("| **Priority** | %s |\n", issue.Priority))
	}
	if issue.Assignee != "" {
		sb.WriteString(fmt.Sprintf("| **Assignee** | %s |\n", issue.Assignee))
	}
	if issue.Reporter != "" {
		sb.WriteString(fmt.Sprintf("| **Reporter** | %s |\n", issue.Reporter))
	}
	if issue.Created != "" {
		sb.WriteString(fmt.Sprintf("| **Created** | %s |\n", formatDate(issue.Created)))
	}
	if issue.Updated != "" {
		sb.WriteString(fmt.Sprintf("| **Updated** | %s |\n", formatDate(issue.Updated)))
	}
	sb.WriteString("\n")

	// Description
	sb.WriteString("## 📝 Description\n\n")
	if issue.Description != "" {
		sb.WriteString(issue.Description)
		sb.WriteString("\n\n")
	} else {
		sb.WriteString("*No description*\n\n")
	}

	// Attachments
	if len(issue.Attachments) > 0 {
		sb.WriteString(fmt.Sprintf("## 📎 Attachments (%d)\n\n", len(issue.Attachments)))
		for i, att := range issue.Attachments {
			sb.WriteString(fmt.Sprintf("%d. **%s** (%s)\n", i+1, att.Filename, formatSize(att.Size)))
			if withImages && isImageFile(att.Filename) {
				// Reference local file
				localPath := filepath.Join("attachments", att.Filename)
				sb.WriteString(fmt.Sprintf("   ![%s](./%s)\n", att.Filename, localPath))
			}
			if att.Author != "" {
				sb.WriteString(fmt.Sprintf("   - Uploaded by: %s\n", att.Author))
			}
			if att.Created != "" {
				sb.WriteString(fmt.Sprintf("   - Created: %s\n", formatDate(att.Created)))
			}
			sb.WriteString("\n")
		}
	}

	// Comments
	if len(issue.Comments) > 0 {
		sb.WriteString(fmt.Sprintf("## 💬 Comments (%d)\n\n", len(issue.Comments)))
		for _, c := range issue.Comments {
			sb.WriteString(fmt.Sprintf("### Comment by %s - %s\n\n", c.Author, formatDate(c.Created)))
			sb.WriteString(fmt.Sprintf("%s\n\n", c.Body))
			sb.WriteString("---\n\n")
		}
	}

	// Links
	sb.WriteString("## 🔗 Links\n\n")
	if f.client != nil {
		url := f.client.GetJiraURL(issue.Key)
		if url != "" {
			sb.WriteString(fmt.Sprintf("- **Jira**: %s\n", url))
		}
	}
	sb.WriteString("\n")

	// Footer
	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("*Exported by qkflow on %s*\n", time.Now().Format("2006-01-02 at 15:04")))
	sb.WriteString("*This is a snapshot. For the latest, visit Jira.*\n")

	return sb.String()
}

// Helper functions

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDate(dateStr string) string {
	// Parse Jira date format (RFC3339)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02 15:04")
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".bmp"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

