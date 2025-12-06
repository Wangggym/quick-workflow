package commands

import (
	"fmt"

	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/spf13/cobra"
)

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Manage Jira issues and status mappings",
	Long: `Work with Jira issues and manage status mappings.

Available commands:
  show    - Display a Jira issue in the terminal
  export  - Export a Jira issue to files (with optional images)
  read    - Intelligently read an issue (recommended for Cursor)
  clean   - Clean up exported files
  list    - List status mappings
  setup   - Setup status mappings
  delete  - Delete status mappings`,
}

var jiraListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Jira status mappings",
	Run: func(cmd *cobra.Command, args []string) {
		statusCache, err := jira.NewStatusCache()
		if err != nil {
			log.Error("Failed to create status cache: %v", err)
			return
		}

		mappings, err := statusCache.ListAllMappings()
		if err != nil {
			log.Error("Failed to list mappings: %v", err)
			return
		}

		if len(mappings) == 0 {
			log.Info("No Jira status mappings configured yet")
			return
		}

		log.Info("\n📋 Jira Status Mappings:")
		log.Info("")
		for _, mapping := range mappings {
			log.Info("Project: %s", mapping.ProjectKey)
			log.Info("  PR Created → %s", mapping.PRCreatedStatus)
			log.Info("  PR Merged  → %s", mapping.PRMergedStatus)
			log.Info("")
		}
	},
}

var jiraSetupCmd = &cobra.Command{
	Use:   "setup [project-key]",
	Short: "Setup or update status mapping for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectKey := args[0]

		jiraClient, err := jira.NewClient()
		if err != nil {
			log.Error("Failed to create Jira client: %v", err)
			return
		}

		mapping, err := setupProjectStatusMapping(jiraClient, projectKey)
		if err != nil {
			log.Error("Failed to setup mapping: %v", err)
			return
		}

		if mapping == nil {
			log.Info("Setup cancelled")
			return
		}

		statusCache, err := jira.NewStatusCache()
		if err != nil {
			log.Error("Failed to create status cache: %v", err)
			return
		}

		if err := statusCache.SaveProjectStatus(mapping); err != nil {
			log.Error("Failed to save mapping: %v", err)
			return
		}

		log.Success("Status mapping saved for project %s!", projectKey)
	},
}

var jiraDeleteCmd = &cobra.Command{
	Use:   "delete [project-key]",
	Short: "Delete status mapping for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectKey := args[0]

		confirm, err := ui.PromptConfirm(fmt.Sprintf("Delete status mapping for %s?", projectKey), false)
		if err != nil || !confirm {
			log.Info("Deletion cancelled")
			return
		}

		statusCache, err := jira.NewStatusCache()
		if err != nil {
			log.Error("Failed to create status cache: %v", err)
			return
		}

		if err := statusCache.DeleteProjectStatus(projectKey); err != nil {
			log.Error("Failed to delete mapping: %v", err)
			return
		}

		log.Success("Status mapping deleted for project %s", projectKey)
	},
}

func init() {
	// Issue commands (new)
	jiraCmd.AddCommand(jiraShowCmd)
	jiraCmd.AddCommand(jiraExportCmd)
	jiraCmd.AddCommand(jiraReadCmd)
	jiraCmd.AddCommand(jiraCleanCmd)

	// Status mapping commands (existing)
	jiraCmd.AddCommand(jiraListCmd)
	jiraCmd.AddCommand(jiraSetupCmd)
	jiraCmd.AddCommand(jiraDeleteCmd)
}
