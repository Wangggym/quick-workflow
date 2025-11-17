package watcher

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier handles desktop notifications
type Notifier struct {
	enabled bool
}

// NewNotifier creates a new Notifier instance
func NewNotifier(enabled bool) *Notifier {
	return &Notifier{
		enabled: enabled,
	}
}

// Notify sends a desktop notification
func (n *Notifier) Notify(title, subtitle, message string) error {
	if !n.enabled {
		return nil
	}

	if runtime.GOOS != "darwin" {
		return nil // Notifications only supported on macOS for now
	}

	// Use osascript to display notification
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s"`, 
		escapeString(message), 
		escapeString(title), 
		escapeString(subtitle))

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// NotifyPRMerged sends a notification when a PR is merged
func (n *Notifier) NotifyPRMerged(prNumber int, title string, tickets []string, updateCount int) error {
	notifTitle := fmt.Sprintf("🎉 PR #%d Merged", prNumber)
	notifSubtitle := truncate(title, 50)
	notifMessage := fmt.Sprintf("%d Jira ticket(s) updated", updateCount)
	
	return n.Notify(notifTitle, notifSubtitle, notifMessage)
}

// NotifyError sends an error notification
func (n *Notifier) NotifyError(title, message string) error {
	return n.Notify("❌ "+title, "", message)
}

// escapeString escapes special characters for AppleScript
func escapeString(s string) string {
	// Escape double quotes and backslashes
	result := ""
	for _, c := range s {
		if c == '"' || c == '\\' {
			result += "\\"
		}
		result += string(c)
	}
	return result
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

