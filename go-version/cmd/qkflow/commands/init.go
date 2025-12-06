package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Wangggym/quick-workflow/internal/config"
	"github.com/Wangggym/quick-workflow/internal/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration for quick-workflow",
	Long: `Initialize configuration by prompting for required credentials and settings.
This will create a configuration file at ~/.config/quick-workflow/config.yaml`,
	Run: runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	log.Info("Welcome to Quick Workflow Setup!")
	showEmptyLine()

	// Reset config cache to ensure we read latest environment variables
	config.Reset()
	// Try to load existing configuration (from YAML or environment variables)
	existingCfg, _ := config.Load()
	cfg := &config.Config{}

	// If we have existing config, use it as defaults for editing
	if existingCfg != nil {
		// Check if configuration is complete and valid
		if existingCfg.Validate() == nil {
			// Configuration is complete, use as defaults for editing
			log.Info("📝 Found existing configuration, editing (press Enter to keep current values)")
			showEmptyLine()
		} else {
			// Configuration exists but incomplete, use as defaults
			log.Info("📋 Found existing configuration (incomplete), using as defaults (press Enter to keep current value)")
			showEmptyLine()
		}
		*cfg = *existingCfg
	}

	// 【第一步】Git 账号配置
	showSectionHeader(1, "Git 账号配置")
	gitAnswers, err := collectGitConfig(cfg)
	if err != nil {
		showError("Failed to collect Git configuration: %v", err)
		return
	}
	cfg.GitHubOwner = gitAnswers.Owner
	cfg.GitHubToken = gitAnswers.Token
	cfg.BranchPrefix = gitAnswers.BranchPrefix

	// 【第二步】Jira 账号配置
	showSectionHeader(2, "Jira 账号配置")
	jiraAnswers, err := collectJiraConfig(cfg)
	if err != nil {
		showError("Failed to collect Jira configuration: %v", err)
		return
	}
	cfg.Email = jiraAnswers.Email
	cfg.JiraServiceAddress = jiraAnswers.ServiceAddress
	cfg.JiraAPIToken = jiraAnswers.APIToken

	// 【第三步】LLM 配置
	showSectionHeader(3, "LLM 配置")
	llmAnswers, err := collectLLMConfig(cfg)
	if err != nil {
		showError("Failed to collect LLM configuration: %v", err)
		return
	}
	if llmAnswers.UseAI {
		cfg.OpenAIKey = llmAnswers.AIKey
	}

	// 【第四步】Log 配置
	showSectionHeader(4, "Log 配置")
	logAnswers, err := collectLogConfig(cfg)
	if err != nil {
		showError("Failed to collect Log configuration: %v", err)
		return
	}
	cfg.LogFilePath = logAnswers.FilePath
	cfg.LogLevel = logAnswers.Level

	// 【第五步】自动更新
	showSectionHeader(5, "自动更新")
	updateAnswers, err := collectUpdateConfig(cfg)
	if err != nil {
		showWarning("Failed to collect update configuration, using default: %v", err)
		updateAnswers = &UpdateConfigAnswers{AutoUpdate: true}
	}
	cfg.AutoUpdate = updateAnswers.AutoUpdate
	if cfg.AutoUpdate {
		showSuccess("Auto-update enabled - qkflow will keep itself up to date")
	} else {
		showInfo("ℹ️", "Auto-update disabled - run 'qkflow update-cli' to update manually")
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		showError("Failed to save configuration: %v", err)
		return
	}

	showSuccess("Configuration saved successfully!")
	showEmptyLine()

	// Show storage location
	showStorageLocation()

	// Show final instructions
	showFinalInstructions()
}

// ==================== Answer Types ====================

// GitConfigAnswers holds answers for Git account configuration
type GitConfigAnswers struct {
	Owner        string
	Token        string
	BranchPrefix string
}

// JiraConfigAnswers holds answers for Jira account configuration
type JiraConfigAnswers struct {
	Email          string
	ServiceAddress string
	APIToken       string
}

// LLMConfigAnswers holds answers for LLM configuration
type LLMConfigAnswers struct {
	UseAI bool
	AIKey string
}

// LogConfigAnswers holds answers for Log configuration
type LogConfigAnswers struct {
	FilePath *string
	Level    *string
}

// UpdateConfigAnswers holds answers for Auto Update configuration
type UpdateConfigAnswers struct {
	AutoUpdate bool
}

// ==================== Output Helper Functions ====================

// showSectionHeader displays a section header with step number and title
func showSectionHeader(step int, title string) {
	log.Info("")
	log.Info("【第%d步】%s", step, title)
	log.Info("──────────────────────────────────────")
	log.Info("")
}

// showConfigList displays a list of configuration items
func showConfigList(title string, items map[string]string) {
	log.Info("%s", title)
	for name, value := range items {
		log.Info("  %s: %s", name, value)
	}
	log.Info("")
}

// showEmptyLine displays an empty line
func showEmptyLine() {
	log.Info("")
}

// showSuccess displays a success message with emoji
func showSuccess(message string, args ...interface{}) {
	if len(args) > 0 {
		log.Success("✅ "+message, args...)
	} else {
		log.Success("✅ " + message)
	}
}

// showError displays an error message with emoji
func showError(message string, args ...interface{}) {
	if len(args) > 0 {
		log.Error("❌ "+message, args...)
	} else {
		log.Error("❌ " + message)
	}
}

// showWarning displays a warning message with emoji
func showWarning(message string, args ...interface{}) {
	if len(args) > 0 {
		log.Warning("⚠️  "+message, args...)
	} else {
		log.Warning("⚠️  " + message)
	}
}

// showInfo displays an info message with emoji
func showInfo(emoji, message string, args ...interface{}) {
	if len(args) > 0 {
		log.Info(emoji+" "+message, args...)
	} else {
		log.Info(emoji + " " + message)
	}
}

// showFinalInstructions displays the final instructions after configuration
func showFinalInstructions() {
	log.Info("")
	log.Info("You can now use the following commands:")
	log.Info("  qkflow pr create   - Create a PR and update Jira")
	log.Info("  qkflow pr merge    - Merge a PR and update Jira")
	log.Info("  qkflow update      - Quick commit and push with PR title")
	log.Info("  qkflow update-cli  - Update qkflow to the latest version")
	log.Info("  qkflow jira list   - List Jira status mappings")
	log.Info("")
}

// showStorageLocation displays the storage location information
func showStorageLocation() {
	location := utils.GetConfigLocation()
	configDir, _ := utils.GetQuickWorkflowConfigDir()
	log.Info("Storage location: %s", location)
	if configDir != "" {
		log.Info("  Config: %s/config.yaml", configDir)
	}
	log.Info("")
}

// ==================== Configuration Collectors ====================

// collectGitConfig collects Git account configuration using survey
func collectGitConfig(existingCfg *config.Config) (*GitConfigAnswers, error) {
	// Try to get GitHub token from gh CLI
	ghToken, err := getGitHubToken()
	hasGhToken := err == nil && ghToken != ""
	if hasGhToken {
		showSuccess("GitHub token obtained from gh CLI")
	} else {
		showWarning("Failed to get GitHub token from gh CLI")
	}

	answers := &GitConfigAnswers{}

	// Build token prompt message
	tokenMessage := "Enter your GitHub personal access token:"
	if existingCfg.GitHubToken != "" {
		// If we have existing token, show keep current message
		tokenMessage = "Enter your GitHub personal access token [current: ****] (press Enter to keep current):"
	} else if hasGhToken {
		// Only show gh CLI message if no existing token configured
		tokenMessage = "Enter your GitHub personal access token [current: ****] (press Enter to use token from gh CLI, or enter a new token):"
	}

	// Token validation: only required if we don't have gh CLI token or existing token
	var tokenValidate survey.Validator
	if !hasGhToken && existingCfg.GitHubToken == "" {
		tokenValidate = survey.Required
	}

	// Custom validator: allow empty if we have existing value, otherwise require
	ownerValidate := func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return fmt.Errorf("invalid input")
		}
		if str == "" && existingCfg.GitHubOwner == "" {
			return fmt.Errorf("github username is required")
		}
		return nil
	}

	qs := []*survey.Question{
		{
			Name: "owner",
			Prompt: &survey.Input{
				Message: func() string {
					if existingCfg.GitHubOwner != "" {
						return fmt.Sprintf("Enter your GitHub username [current: %s] (press Enter to keep):", existingCfg.GitHubOwner)
					}
					return "Enter your GitHub username:"
				}(),
				// Don't set Default to avoid showing (default_value), handle empty input manually
			},
			Validate: ownerValidate,
		},
		{
			Name: "token",
			Prompt: &survey.Password{
				Message: tokenMessage,
			},
			Validate: tokenValidate,
		},
		{
			Name: "branchPrefix",
			Prompt: &survey.Input{
				Message: func() string {
					if existingCfg.BranchPrefix != "" {
						return fmt.Sprintf("Enter branch prefix [current: %s] (optional, press Enter to keep):", existingCfg.BranchPrefix)
					}
					return "Enter branch prefix (optional, e.g., 'feature' or 'username'):"
				}(),
				// Don't set Default to avoid showing (default_value), handle empty input manually
			},
		},
	}

	if err := survey.Ask(qs, answers); err != nil {
		return nil, err
	}

	// Handle owner: if user entered nothing and we have existing value, use it
	if answers.Owner == "" && existingCfg.GitHubOwner != "" {
		answers.Owner = existingCfg.GitHubOwner
	}

	// Handle token: if user entered nothing and we have gh CLI token or existing token, use it
	if answers.Token == "" {
		if hasGhToken {
			answers.Token = ghToken
		} else if existingCfg.GitHubToken != "" {
			answers.Token = existingCfg.GitHubToken
		}
	}

	// Handle branch prefix: if empty, use existing value; trim trailing slash
	if answers.BranchPrefix == "" && existingCfg.BranchPrefix != "" {
		answers.BranchPrefix = existingCfg.BranchPrefix
	}
	answers.BranchPrefix = strings.TrimRight(answers.BranchPrefix, "/")

	return answers, nil
}

// collectJiraConfig collects Jira account configuration using survey
func collectJiraConfig(existingCfg *config.Config) (*JiraConfigAnswers, error) {
	answers := &JiraConfigAnswers{}

	// Custom validators: allow empty if we have existing value, otherwise require
	emailValidate := func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return fmt.Errorf("invalid input")
		}
		if str == "" && existingCfg.Email == "" {
			return fmt.Errorf("jira email address is required")
		}
		return nil
	}

	serviceAddressValidate := func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return fmt.Errorf("invalid input")
		}
		if str == "" && existingCfg.JiraServiceAddress == "" {
			return fmt.Errorf("jira service address is required")
		}
		return nil
	}

	qs := []*survey.Question{
		{
			Name: "email",
			Prompt: &survey.Input{
				Message: func() string {
					if existingCfg.Email != "" {
						return fmt.Sprintf("Enter your Jira email address [current: %s] (press Enter to keep):", existingCfg.Email)
					}
					return "Enter your Jira email address:"
				}(),
				// Don't set Default to avoid showing (default_value), handle empty input manually
			},
			Validate: emailValidate,
		},
		{
			Name: "serviceAddress",
			Prompt: &survey.Input{
				Message: func() string {
					if existingCfg.JiraServiceAddress != "" {
						return fmt.Sprintf("Enter your Jira service address [current: %s] (press Enter to keep):", existingCfg.JiraServiceAddress)
					}
					return "Enter your Jira service address (e.g., https://your-domain.atlassian.net):"
				}(),
				// Don't set Default to avoid showing (default_value), handle empty input manually
			},
			Validate: serviceAddressValidate,
		},
		{
			Name: "apiToken",
			Prompt: &survey.Password{
				Message: func() string {
					if existingCfg.JiraAPIToken != "" {
						return "Enter your Jira API token [current: ***] (press Enter to keep current):"
					}
					return "Enter your Jira API token:"
				}(),
			},
			// Token validation: only required if we don't have existing token
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return fmt.Errorf("invalid input")
				}
				if str == "" && existingCfg.JiraAPIToken == "" {
					return fmt.Errorf("jira API token is required")
				}
				return nil
			},
		},
	}

	showInfo("ℹ️", "Get your Jira API token from: https://id.atlassian.com/manage-profile/security/api-tokens")

	if err := survey.Ask(qs, answers); err != nil {
		return nil, err
	}

	// Handle email: if user entered nothing and we have existing value, use it
	if answers.Email == "" && existingCfg.Email != "" {
		answers.Email = existingCfg.Email
	}

	// Handle service address: if user entered nothing and we have existing value, use it
	if answers.ServiceAddress == "" && existingCfg.JiraServiceAddress != "" {
		answers.ServiceAddress = existingCfg.JiraServiceAddress
	}

	// Handle API token: if user entered nothing and we have existing token, use it
	if answers.APIToken == "" && existingCfg.JiraAPIToken != "" {
		answers.APIToken = existingCfg.JiraAPIToken
	}

	// Handle service address: trim trailing slash
	answers.ServiceAddress = strings.TrimRight(answers.ServiceAddress, "/")

	return answers, nil
}

// collectLLMConfig collects LLM configuration using survey
func collectLLMConfig(existingCfg *config.Config) (*LLMConfigAnswers, error) {
	hasAIKey := existingCfg.OpenAIKey != "" || existingCfg.DeepSeekKey != ""

	answers := &LLMConfigAnswers{
		UseAI: hasAIKey,
	}

	// First, ask if user wants to configure AI
	useAIPrompt := &survey.Confirm{
		Message: "Do you want to configure AI features (OpenAI/DeepSeek)?",
		Default: hasAIKey,
	}
	if err := survey.AskOne(useAIPrompt, &answers.UseAI); err != nil {
		return nil, err
	}

	// If user wants to configure AI, ask for the key
	if answers.UseAI {
		aiKeyMessage := "Enter OpenAI or DeepSeek API key:"
		if existingCfg.OpenAIKey != "" || existingCfg.DeepSeekKey != "" {
			aiKeyMessage = "Enter OpenAI or DeepSeek API key [current: ***] (press Enter to keep current):"
		}

		aiKeyPrompt := &survey.Password{
			Message: aiKeyMessage,
		}

		var aiKey string
		if err := survey.AskOne(aiKeyPrompt, &aiKey); err != nil {
			return nil, err
		}

		// Handle AI key: if user entered nothing and we have existing key, use it
		if aiKey == "" {
			// Prefer OpenAIKey, fallback to DeepSeekKey
			if existingCfg.OpenAIKey != "" {
				answers.AIKey = existingCfg.OpenAIKey
			} else if existingCfg.DeepSeekKey != "" {
				answers.AIKey = existingCfg.DeepSeekKey
			}
		} else {
			answers.AIKey = aiKey
		}
	} else {
		// User doesn't want to configure AI, clear the key
		answers.AIKey = ""
	}

	return answers, nil
}

// collectLogConfig collects Log configuration using survey
func collectLogConfig(existingCfg *config.Config) (*LogConfigAnswers, error) {
	answers := &LogConfigAnswers{
		FilePath: existingCfg.LogFilePath,
		Level:    existingCfg.LogLevel,
	}

	var filePathDefault string
	var levelDefault string
	if answers.FilePath != nil {
		filePathDefault = *answers.FilePath
	}
	if answers.Level != nil {
		levelDefault = *answers.Level
	}

	qs := []*survey.Question{
		{
			Name: "filePath",
			Prompt: &survey.Input{
				Message: "Enter log file path (optional, press Enter to skip):",
				Default: filePathDefault,
			},
		},
		{
			Name: "level",
			Prompt: &survey.Select{
				Message: "Select log level (optional):",
				Options: []string{"skip", "debug", "info", "warn", "error"},
				Default: func() string {
					if levelDefault == "" {
						return "skip"
					}
					return levelDefault
				}(),
			},
		},
	}

	var result struct {
		FilePath string
		Level    string
	}

	if err := survey.Ask(qs, &result); err != nil {
		return nil, err
	}

	// Convert to pointers: empty string becomes nil, non-empty becomes *string
	if result.FilePath != "" {
		answers.FilePath = &result.FilePath
	} else {
		answers.FilePath = nil
	}

	// "skip" means nil, otherwise use the selected value
	if result.Level != "" && result.Level != "skip" {
		answers.Level = &result.Level
	} else {
		answers.Level = nil
	}

	return answers, nil
}

// collectUpdateConfig collects Auto Update configuration using survey
func collectUpdateConfig(existingCfg *config.Config) (*UpdateConfigAnswers, error) {
	defaultValue := true
	if existingCfg != nil {
		defaultValue = existingCfg.AutoUpdate
	}

	answers := &UpdateConfigAnswers{
		AutoUpdate: defaultValue,
	}

	qs := []*survey.Question{
		{
			Name: "autoUpdate",
			Prompt: &survey.Confirm{
				Message: "Enable automatic updates? (recommended)",
				Default: defaultValue,
			},
		},
	}

	if err := survey.Ask(qs, answers); err != nil {
		// On error, default to true
		answers.AutoUpdate = true
	}

	return answers, nil
}

// getGitHubToken gets GitHub token from gh CLI
func getGitHubToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
