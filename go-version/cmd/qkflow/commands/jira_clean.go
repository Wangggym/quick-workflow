package commands

import (
	"fmt"

	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/spf13/cobra"
)

var (
	cleanAll    bool
	cleanForce  bool
	cleanDryRun bool
)

var jiraCleanCmd = &cobra.Command{
	Use:   "clean [ISSUE-KEY]",
	Short: "Clean up exported Jira files",
	Long: `Delete exported Jira files to free up disk space.

By default, prompts for confirmation before deleting.
Use --force to skip confirmation.
Use --dry-run to see what would be deleted without actually deleting.

Examples:
  qkflow jira clean NA-9245           # Clean specific issue
  qkflow jira clean --all             # Clean all exports
  qkflow jira clean --all --force     # Clean all without confirmation
  qkflow jira clean NA-9245 --dry-run # Preview what would be deleted`,
	Args: func(cmd *cobra.Command, args []string) error {
		if cleanAll && len(args) > 0 {
			return fmt.Errorf("cannot specify issue key with --all flag")
		}
		if !cleanAll && len(args) == 0 {
			return fmt.Errorf("issue key is required (or use --all to clean all exports)")
		}
		if !cleanAll && len(args) != 1 {
			return fmt.Errorf("exactly one issue key is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Create cleaner
		cleaner := jira.NewCleaner()

		// Prepare options
		opts := jira.CleanOptions{
			All:    cleanAll,
			DryRun: cleanDryRun,
		}

		if !cleanAll {
			issueKey := args[0]
			if !jira.ValidateIssueKey(issueKey) {
				log.Error("Invalid issue key: %s", issueKey)
				return
			}
			opts.IssueKey = issueKey
		}

		// Preview what will be cleaned
		if cleanAll && !cleanForce && !cleanDryRun {
			// List all exports
			exports, err := cleaner.ListExports()
			if err != nil {
				log.Error("Failed to list exports: %v", err)
				return
			}

			if len(exports) == 0 {
				log.Info("No exports found.")
				return
			}

			log.Info("🗑️  The following will be deleted:")
			totalSize := int64(0)
			for _, exp := range exports {
				log.Info("  • %s (%s, %d files)", exp.IssueKey, formatSize(exp.Size), exp.FileCount)
				totalSize += exp.Size
			}
			log.Info("\nTotal: %d exports, %s\n", len(exports), formatSize(totalSize))

			// Ask for confirmation
			confirm, err := ui.PromptConfirm("Are you sure you want to delete all exports?", false)
			if err != nil || !confirm {
				log.Info("Cancelled.")
				return
			}
		} else if !cleanAll && !cleanForce && !cleanDryRun {
			// Single issue: show preview and confirm
			// Try to get size first
			testResult, err := cleaner.Clean(jira.CleanOptions{
				IssueKey: opts.IssueKey,
				DryRun:   true,
			})
			if err != nil {
				log.Error("Failed to check export: %v", err)
				return
			}

			if len(testResult) > 0 && testResult[0].Error == nil {
				r := testResult[0]
				log.Info("🗑️  The following will be deleted:")
				log.Info("   %s (%s, %d files)", r.IssueKey, formatSize(r.Size), r.FileCount)

				confirm, err := ui.PromptConfirm(fmt.Sprintf("Delete export for %s?", opts.IssueKey), false)
				if err != nil || !confirm {
					log.Info("Cancelled.")
					return
				}
			}
		}

		// Perform cleaning
		results, err := cleaner.Clean(opts)
		if err != nil {
			log.Error("Failed to clean: %v", err)
			return
		}

		// Display results
		log.Info("")
		if cleanDryRun {
			log.Info("🔍 Dry run - no files were deleted:")
		}

		totalSize := int64(0)
		successCount := 0
		for _, result := range results {
			fmt.Print(jira.FormatCleanResult(result, cleanDryRun))
			if result.Error == nil {
				totalSize += result.Size
				if result.Deleted {
					successCount++
				}
			}
		}

		log.Info("")
		if cleanDryRun {
			log.Info("Would free: %s", formatSize(totalSize))
		} else {
			if successCount > 0 {
				log.Success("Cleaned %d export(s), freed %s", successCount, formatSize(totalSize))
			}
		}
	},
}

func init() {
	jiraCleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean all exported issues")
	jiraCleanCmd.Flags().BoolVarP(&cleanForce, "force", "f", false, "Skip confirmation prompt")
	jiraCleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
}

// Helper function
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
