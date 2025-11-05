package commands

import (
	"fmt"

	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/Wangggym/quick-workflow/pkg/config"
	"github.com/spf13/cobra"
)

var (
	// Version is the application version
	Version = "dev"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "qk",
	Short: "Quick workflow tool for GitHub and Jira",
	Long: `qk is a CLI tool to streamline your GitHub and Jira workflow.
It automates common tasks like creating PRs, updating Jira status, and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 对于某些命令不需要检查配置
		skipConfigCheck := []string{"init", "version", "help"}
		for _, skip := range skipConfigCheck {
			if cmd.Name() == skip || cmd.Parent().Name() == skip {
				return
			}
		}

		// 检查配置
		if _, err := config.Load(); err != nil {
			ui.Error(fmt.Sprintf("Failed to load config: %v", err))
			ui.Warning("Please run 'qk init' to configure the tool")
			return
		}

		if !config.IsConfigured() {
			ui.Warning("Configuration incomplete. Please run 'qk init' to complete setup")
		}
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(prCmd)
	rootCmd.AddCommand(configCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("qk version %s (built: %s)\n", Version, BuildTime)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		if cfg == nil {
			ui.Error("No configuration found. Run 'qk init' first.")
			return
		}

		fmt.Println("Current configuration:")
		fmt.Printf("  Email: %s\n", cfg.Email)
		fmt.Printf("  Jira Service: %s\n", cfg.JiraServiceAddress)
		fmt.Printf("  GitHub Token: %s\n", maskToken(cfg.GitHubToken))
		fmt.Printf("  Jira API Token: %s\n", maskToken(cfg.JiraAPIToken))
		if cfg.BranchPrefix != "" {
			fmt.Printf("  Branch Prefix: %s\n", cfg.BranchPrefix)
		}
	},
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

