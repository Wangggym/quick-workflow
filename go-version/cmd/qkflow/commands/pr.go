package commands

import (
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Pull request operations",
	Long:  `Create, merge, and manage pull requests with automatic Jira integration.`,
}

func init() {
	prCmd.AddCommand(prCreateCmd)
	prCmd.AddCommand(prMergeCmd)
	prCmd.AddCommand(prApproveCmd)
}

