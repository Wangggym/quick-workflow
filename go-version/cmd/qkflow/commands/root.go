package commands

import (
	"github.com/Wangggym/quick-workflow/internal/config"
	"github.com/Wangggym/quick-workflow/internal/logger"
	"github.com/Wangggym/quick-workflow/internal/updater"
	"github.com/Wangggym/quick-workflow/internal/utils"
	"github.com/spf13/cobra"
)

// Global logger instance for all commands
var log *logger.Logger

func init() {
	// Initialize global logger for commands
	// Level will be automatically loaded from environment variable or default value
	log, _ = logger.NewLogger(&logger.LoggerOptions{
		Type: logger.LoggerTypeUI,
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
	})
}

var (
	// Version is the application version
	Version = "dev"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "qkflow",
	Short:   "Quick workflow tool for GitHub and Jira",
	Version: Version,
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
			log.Error("Failed to load config: %v", err)
			log.Warning("Please run 'qkflow init' to configure the tool")
			return
		}

		if !config.IsConfigured() {
			log.Warning("Configuration incomplete. Please run 'qkflow init' to complete setup")
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

// SetVersion sets the version string for the root command
func SetVersion(version string) {
	rootCmd.Version = version
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(prCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(jiraCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(watchCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("qkflow version %s (built: %s)", Version, BuildTime)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		if cfg == nil {
			log.Error("No configuration found. Run 'qkflow init' first.")
			return
		}

		log.Info("Current configuration:")
		log.Info("")

		// Show storage location
		location := utils.GetConfigLocation()
		configDir, _ := utils.GetQuickWorkflowConfigDir()
		jiraDir, _ := utils.GetConfigDir()
		log.Info("💾 Storage:")
		log.Info("  Location: %s", location)
		if configDir != "" {
			log.Info("  Config: %s/config.yaml", configDir)
		}
		if jiraDir != "" {
			log.Info("  Jira Status: %s/jira-status.json", jiraDir)
		}

		log.Info("")
		log.Info("📧 Basic:")
		log.Info("  Email: %s", cfg.Email)
		if cfg.BranchPrefix != "" {
			log.Info("  Branch Prefix: %s", cfg.BranchPrefix)
		}

		log.Info("")
		log.Info("🐙 GitHub:")
		log.Info("  Token: %s", maskToken(cfg.GitHubToken))
		if cfg.GitHubOwner != "" {
			log.Info("  Owner: %s", cfg.GitHubOwner)
		}
		if cfg.GitHubRepo != "" {
			log.Info("  Repo: %s", cfg.GitHubRepo)
		}

		log.Info("")
		log.Info("📋 Jira:")
		log.Info("  Service: %s", cfg.JiraServiceAddress)
		log.Info("  API Token: %s", maskToken(cfg.JiraAPIToken))

		log.Info("")
		log.Info("🔄 Auto Update:")
		if cfg.AutoUpdate {
			log.Info("  Status: ✅ Enabled (checks every 24h)")
		} else {
			log.Info("  Status: ❌ Disabled (run 'qkflow update-cli' to update manually)")
		}

		log.Info("")
		log.Info("🤖 AI (optional):")

		// Determine which AI service is active
		hasDeepSeek := cfg.DeepSeekKey != ""
		hasOpenAI := cfg.OpenAIKey != ""

		if hasDeepSeek {
			log.Info("  DeepSeek Key: %s ✅ (Active)", maskToken(cfg.DeepSeekKey))
		} else {
			log.Info("  DeepSeek Key: not configured")
		}

		if hasOpenAI {
			if hasDeepSeek {
				log.Info("  OpenAI Key: %s ✅ (Backup)", maskToken(cfg.OpenAIKey))
			} else {
				log.Info("  OpenAI Key: %s ✅ (Active)", maskToken(cfg.OpenAIKey))
			}
		} else {
			log.Info("  OpenAI Key: not configured")
		}

		if cfg.OpenAIProxyURL != "" {
			log.Info("  OpenAI Proxy URL: %s", cfg.OpenAIProxyURL)
		}

		if !hasDeepSeek && !hasOpenAI {
			log.Info("")
			log.Info("  💡 Tip: Configure AI for automatic PR title/description generation")
		}
	},
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
