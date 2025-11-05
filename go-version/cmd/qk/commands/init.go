package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Wangggym/quick-workflow/internal/ui"
	"github.com/Wangggym/quick-workflow/pkg/config"
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
	ui.Info("Welcome to Quick Workflow Setup!")
	fmt.Println()

	cfg := &config.Config{}

	// Email
	email, err := ui.PromptInput("Enter your email address:", true)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get email: %v", err))
		return
	}
	cfg.Email = email

	// GitHub Token
	ui.Info("Getting GitHub token from gh CLI...")
	ghToken, err := getGitHubToken()
	if err != nil {
		ui.Warning("Failed to get GitHub token from gh CLI")
		ghToken, err = ui.PromptPassword("Enter your GitHub personal access token:")
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to get GitHub token: %v", err))
			return
		}
	} else {
		ui.Success("GitHub token obtained from gh CLI")
	}
	cfg.GitHubToken = ghToken

	// Jira Service Address
	jiraAddr, err := ui.PromptInput("Enter your Jira service address (e.g., https://your-domain.atlassian.net):", true)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get Jira address: %v", err))
		return
	}
	cfg.JiraServiceAddress = strings.TrimRight(jiraAddr, "/")

	// Jira API Token
	ui.Info("Get your Jira API token from: https://id.atlassian.com/manage-profile/security/api-tokens")
	jiraToken, err := ui.PromptPassword("Enter your Jira API token:")
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get Jira token: %v", err))
		return
	}
	cfg.JiraAPIToken = jiraToken

	// Branch Prefix (optional)
	branchPrefix, err := ui.PromptInput("Enter branch prefix (optional, e.g., 'feature' or 'username'):", false)
	if err != nil {
		ui.Warning("Skipping branch prefix")
	} else {
		cfg.BranchPrefix = strings.TrimRight(branchPrefix, "/")
	}

	// OpenAI Key (optional)
	useAI, err := ui.PromptConfirm("Do you want to configure AI features (OpenAI/DeepSeek)?", false)
	if err == nil && useAI {
		aiKey, err := ui.PromptPassword("Enter OpenAI or DeepSeek API key:")
		if err == nil {
			cfg.OpenAIKey = aiKey
		}
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		ui.Error(fmt.Sprintf("Failed to save configuration: %v", err))
		return
	}

	ui.Success("Configuration saved successfully!")
	fmt.Println()
	ui.Info("You can now use the following commands:")
	fmt.Println("  qk pr create  - Create a PR and update Jira")
	fmt.Println("  qk pr merge   - Merge a PR and update Jira")
	fmt.Println()
}

func getGitHubToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

