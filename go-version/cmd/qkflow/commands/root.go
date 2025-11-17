package commands

import (
	"fmt"

	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/Wangggym/quick-workflow/internal/utils"
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
	Use:   "qkflow",
	Short: "Quick workflow tool for GitHub and Jira",
	Long: `qkflow is a CLI tool to streamline your GitHub and Jira workflow.
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
			ui.Warning("Please run 'qkflow init' to configure the tool")
			return
		}

		if !config.IsConfigured() {
			ui.Warning("Configuration incomplete. Please run 'qkflow init' to complete setup")
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
	rootCmd.AddCommand(jiraCmd)
	rootCmd.AddCommand(updateCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("qkflow version %s (built: %s)\n", Version, BuildTime)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		if cfg == nil {
			ui.Error("No configuration found. Run 'qkflow init' first.")
			return
		}

		fmt.Println("Current configuration:")
		fmt.Println()
		
		// Show storage location
		location := utils.GetConfigLocation()
		configDir, _ := utils.GetQuickWorkflowConfigDir()
		jiraDir, _ := utils.GetConfigDir()
		fmt.Println("💾 Storage:")
		fmt.Printf("  Location: %s\n", location)
		if configDir != "" {
			fmt.Printf("  Config: %s/config.yaml\n", configDir)
		}
		if jiraDir != "" {
			fmt.Printf("  Jira Status: %s/jira-status.json\n", jiraDir)
		}
		
		fmt.Println()
		fmt.Println("📧 Basic:")
		fmt.Printf("  Email: %s\n", cfg.Email)
		if cfg.BranchPrefix != "" {
			fmt.Printf("  Branch Prefix: %s\n", cfg.BranchPrefix)
		}
		
		fmt.Println()
		fmt.Println("🐙 GitHub:")
		fmt.Printf("  Token: %s\n", maskToken(cfg.GitHubToken))
		
		fmt.Println()
		fmt.Println("📋 Jira:")
		fmt.Printf("  Service: %s\n", cfg.JiraServiceAddress)
		fmt.Printf("  API Token: %s\n", maskToken(cfg.JiraAPIToken))
		
		fmt.Println()
		fmt.Println("🤖 AI (optional):")
		if cfg.DeepSeekKey != "" {
			fmt.Printf("  DeepSeek Key: %s ✅\n", maskToken(cfg.DeepSeekKey))
		} else {
			fmt.Printf("  DeepSeek Key: not configured\n")
		}
		if cfg.OpenAIKey != "" {
			fmt.Printf("  OpenAI Key: %s ✅\n", maskToken(cfg.OpenAIKey))
		} else {
			fmt.Printf("  OpenAI Key: not configured\n")
		}
		if cfg.OpenAIProxyURL != "" {
			fmt.Printf("  OpenAI Proxy URL: %s\n", cfg.OpenAIProxyURL)
		}
	},
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

