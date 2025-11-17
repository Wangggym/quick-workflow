package commands

import (
	"fmt"

	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/Wangggym/quick-workflow/internal/updater"
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
		// 对于某些命令不需要检查配置和更新
		skipConfigCheck := []string{"init", "version", "help", "update-cli"}
		for _, skip := range skipConfigCheck {
			if cmd.Name() == skip || cmd.Parent().Name() == skip {
				return
			}
		}

		// 检查配置
		cfg, err := config.Load()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to load config: %v", err))
			ui.Warning("Please run 'qkflow init' to configure the tool")
			return
		}

		if !config.IsConfigured() {
			ui.Warning("Configuration incomplete. Please run 'qkflow init' to complete setup")
		}

		// 检查更新（后台静默执行，不阻塞主流程）
		go func() {
			if err := updater.CheckAndUpdate(Version, cfg.AutoUpdate); err != nil {
				// 静默失败，不影响主流程
			}
		}()
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
		if cfg.GitHubOwner != "" {
			fmt.Printf("  Owner: %s\n", cfg.GitHubOwner)
		}
		if cfg.GitHubRepo != "" {
			fmt.Printf("  Repo: %s\n", cfg.GitHubRepo)
		}
		
		fmt.Println()
		fmt.Println("📋 Jira:")
		fmt.Printf("  Service: %s\n", cfg.JiraServiceAddress)
		fmt.Printf("  API Token: %s\n", maskToken(cfg.JiraAPIToken))
		
		fmt.Println()
		fmt.Println("🔄 Auto Update:")
		if cfg.AutoUpdate {
			fmt.Printf("  Status: ✅ Enabled (checks every 24h)\n")
		} else {
			fmt.Printf("  Status: ❌ Disabled (run 'qkflow update-cli' to update manually)\n")
		}
		
		fmt.Println()
		fmt.Println("🤖 AI (optional):")
		
		// Determine which AI service is active
		hasDeepSeek := cfg.DeepSeekKey != ""
		hasOpenAI := cfg.OpenAIKey != ""
		
		if hasDeepSeek {
			fmt.Printf("  DeepSeek Key: %s ✅ (Active)\n", maskToken(cfg.DeepSeekKey))
		} else {
			fmt.Printf("  DeepSeek Key: not configured\n")
		}
		
		if hasOpenAI {
			if hasDeepSeek {
				fmt.Printf("  OpenAI Key: %s ✅ (Backup)\n", maskToken(cfg.OpenAIKey))
			} else {
				fmt.Printf("  OpenAI Key: %s ✅ (Active)\n", maskToken(cfg.OpenAIKey))
			}
		} else {
			fmt.Printf("  OpenAI Key: not configured\n")
		}
		
		if cfg.OpenAIProxyURL != "" {
			fmt.Printf("  OpenAI Proxy URL: %s\n", cfg.OpenAIProxyURL)
		}
		
		if !hasDeepSeek && !hasOpenAI {
			fmt.Println()
			fmt.Println("  💡 Tip: Configure AI for automatic PR title/description generation")
		}
	},
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

