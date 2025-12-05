package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Wangggym/quick-workflow/internal/config"
	"github.com/Wangggym/quick-workflow/internal/ui"
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
	log.Info("")

	// Reset config cache to ensure we read latest environment variables
	config.Reset()
	// Try to load existing configuration (from YAML or environment variables)
	existingCfg, _ := config.Load()
	cfg := &config.Config{}

	// If we have existing config, check if it's complete and valid
	if existingCfg != nil {
		// Check if configuration is complete and valid
		if existingCfg.Validate() == nil {
			// Configuration is complete, ask if user wants to use it directly
			log.Info("📋 Found existing complete configuration:")
			log.Info("  Email: %s", existingCfg.Email)
			log.Info("  GitHub Owner: %s", existingCfg.GitHubOwner)
			log.Info("  GitHub Repo: %s", existingCfg.GitHubRepo)
			log.Info("  Jira Service: %s", existingCfg.JiraServiceAddress)
			log.Info("")

			useExisting, err := ui.PromptConfirm("Use this configuration? (press Enter to use, or 'n' to edit)", true)
			if err == nil && useExisting {
				// User wants to use existing config, just save it (in case it was from env vars)
				if err := config.Save(existingCfg); err != nil {
					log.Error("Failed to save configuration: %v", err)
					return
				}
				log.Success("Configuration saved successfully!")
				log.Info("")
				return
			}
			// User wants to edit, continue with prompts using existing as defaults
			log.Info("📝 Editing configuration (press Enter to keep current values)")
			log.Info("")
		} else {
			// Configuration exists but incomplete, use as defaults
			log.Info("📋 Found existing configuration (incomplete), using as defaults (press Enter to keep current value)")
			log.Info("")
		}
		*cfg = *existingCfg
	}

	// Email
	emailPrompt := "Enter your email address:"
	if cfg.Email != "" {
		emailPrompt = fmt.Sprintf("Enter your email address [%s]:", cfg.Email)
	}
	email, err := ui.PromptInputWithDefault(emailPrompt, cfg.Email, true)
	if err != nil {
		log.Error("Failed to get email: %v", err)
		return
	}
	cfg.Email = email

	// GitHub Token
	log.Info("Getting GitHub token from gh CLI...")
	ghToken, err := getGitHubToken()
	if err != nil {
		log.Warning("Failed to get GitHub token from gh CLI")
		// Try to use existing token as default
		ghTokenPrompt := "Enter your GitHub personal access token:"
		if cfg.GitHubToken != "" {
			ghTokenPrompt = "Enter your GitHub personal access token (press Enter to keep current):"
		}
		ghToken, err = ui.PromptPasswordWithDefault(ghTokenPrompt, cfg.GitHubToken)
		if err != nil {
			log.Error("Failed to get GitHub token: %v", err)
			return
		}
	} else {
		log.Success("GitHub token obtained from gh CLI")
	}
	cfg.GitHubToken = ghToken

	// GitHub Owner
	githubOwnerPrompt := "Enter your GitHub username or organization:"
	if cfg.GitHubOwner != "" {
		githubOwnerPrompt = fmt.Sprintf("Enter your GitHub username or organization [%s]:", cfg.GitHubOwner)
	}
	githubOwner, err := ui.PromptInputWithDefault(githubOwnerPrompt, cfg.GitHubOwner, true)
	if err != nil {
		log.Error("Failed to get GitHub owner: %v", err)
		return
	}
	cfg.GitHubOwner = githubOwner

	// GitHub Repo
	githubRepoPrompt := "Enter your GitHub repository name:"
	if cfg.GitHubRepo != "" {
		githubRepoPrompt = fmt.Sprintf("Enter your GitHub repository name [%s]:", cfg.GitHubRepo)
	}
	githubRepo, err := ui.PromptInputWithDefault(githubRepoPrompt, cfg.GitHubRepo, true)
	if err != nil {
		log.Error("Failed to get GitHub repo: %v", err)
		return
	}
	cfg.GitHubRepo = githubRepo

	// Jira Service Address
	jiraAddrPrompt := "Enter your Jira service address (e.g., https://your-domain.atlassian.net):"
	if cfg.JiraServiceAddress != "" {
		jiraAddrPrompt = fmt.Sprintf("Enter your Jira service address [%s]:", cfg.JiraServiceAddress)
	}
	jiraAddr, err := ui.PromptInputWithDefault(jiraAddrPrompt, cfg.JiraServiceAddress, true)
	if err != nil {
		log.Error("Failed to get Jira address: %v", err)
		return
	}
	cfg.JiraServiceAddress = strings.TrimRight(jiraAddr, "/")

	// Jira API Token
	log.Info("Get your Jira API token from: https://id.atlassian.com/manage-profile/security/api-tokens")
	jiraTokenPrompt := "Enter your Jira API token:"
	if cfg.JiraAPIToken != "" {
		jiraTokenPrompt = "Enter your Jira API token (press Enter to keep current):"
	}
	jiraToken, err := ui.PromptPasswordWithDefault(jiraTokenPrompt, cfg.JiraAPIToken)
	if err != nil {
		log.Error("Failed to get Jira token: %v", err)
		return
	}
	cfg.JiraAPIToken = jiraToken

	// Branch Prefix (optional)
	branchPrefixPrompt := "Enter branch prefix (optional, e.g., 'feature' or 'username'):"
	if cfg.BranchPrefix != "" {
		branchPrefixPrompt = fmt.Sprintf("Enter branch prefix [%s] (optional, press Enter to keep):", cfg.BranchPrefix)
	}
	branchPrefix, err := ui.PromptInputWithDefault(branchPrefixPrompt, cfg.BranchPrefix, false)
	if err != nil {
		log.Warning("Skipping branch prefix")
	} else {
		cfg.BranchPrefix = strings.TrimRight(branchPrefix, "/")
	}

	// OpenAI Key (optional)
	hasAIKey := cfg.OpenAIKey != "" || cfg.DeepSeekKey != ""
	useAI, err := ui.PromptConfirm("Do you want to configure AI features (OpenAI/DeepSeek)?", hasAIKey)
	if err == nil && useAI {
		aiKeyPrompt := "Enter OpenAI or DeepSeek API key:"
		if cfg.OpenAIKey != "" {
			aiKeyPrompt = "Enter OpenAI or DeepSeek API key (press Enter to keep current):"
		}
		aiKey, err := ui.PromptPasswordWithDefault(aiKeyPrompt, cfg.OpenAIKey)
		if err == nil && aiKey != "" {
			cfg.OpenAIKey = aiKey
		}
	}

	// Auto Update (default: true)
	autoUpdate, err := ui.PromptConfirm("Enable automatic updates? (recommended)", true)
	if err == nil {
		cfg.AutoUpdate = autoUpdate
	} else {
		cfg.AutoUpdate = true // Default to true
	}
	if cfg.AutoUpdate {
		log.Success("✅ Auto-update enabled - qkflow will keep itself up to date")
	} else {
		log.Info("ℹ️  Auto-update disabled - run 'qkflow update-cli' to update manually")
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		log.Error("Failed to save configuration: %v", err)
		return
	}

	log.Success("Configuration saved successfully!")
	log.Info("")

	// Show storage location
	location := utils.GetConfigLocation()
	configDir, _ := utils.GetQuickWorkflowConfigDir()
	log.Info("Storage location: %s", location)
	if configDir != "" {
		log.Info("  📁 Config: %s/config.yaml", configDir)
	}
	log.Info("")

	log.Info("You can now use the following commands:")
	log.Info("  qkflow pr create   - Create a PR and update Jira")
	log.Info("  qkflow pr merge    - Merge a PR and update Jira")
	log.Info("  qkflow update      - Quick commit and push with PR title")
	log.Info("  qkflow update-cli  - Update qkflow to the latest version")
	log.Info("  qkflow jira list   - List Jira status mappings")
	log.Info("")
}

func getGitHubToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
