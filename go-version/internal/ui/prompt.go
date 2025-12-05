package ui

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

var (
	// 颜色定义（用于交互式 prompt）
	Green   = color.New(color.FgGreen).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
)

// PromptInput prompts for a text input
func PromptInput(message string, required bool) (string, error) {
	return PromptInputWithDefault(message, "", required)
}

// PromptInputWithDefault prompts for a text input with a default value
func PromptInputWithDefault(message string, defaultValue string, required bool) (string, error) {
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}

	var result string
	var opts []survey.AskOpt
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	if err := survey.AskOne(prompt, &result, opts...); err != nil {
		return "", err
	}

	return result, nil
}

// PromptPassword prompts for a password input
func PromptPassword(message string) (string, error) {
	return PromptPasswordWithDefault(message, "")
}

// PromptPasswordWithDefault prompts for a password input with optional default value hint
// Note: For security, we don't show the actual default value, but can use it if user presses Enter
func PromptPasswordWithDefault(message string, defaultValue string) (string, error) {
	prompt := &survey.Password{
		Message: message,
	}

	var result string
	var opts []survey.AskOpt

	// If default value exists, make it optional (user can press Enter to use default)
	if defaultValue != "" {
		// We'll handle the default value logic after getting input
	} else {
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	if err := survey.AskOne(prompt, &result, opts...); err != nil {
		return "", err
	}

	// If user entered nothing and we have a default, use it
	if result == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return result, nil
}

// PromptConfirm prompts for a yes/no confirmation
func PromptConfirm(message string, defaultValue bool) (bool, error) {
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}

	var result bool
	if err := survey.AskOne(prompt, &result); err != nil {
		return false, err
	}

	return result, nil
}

// PromptOptional prompts for an optional action with space to select, enter to skip
// This provides better UX: Space = Yes, Enter = Skip (default)
func PromptOptional(message string) (bool, error) {
	options := []string{
		"⏭️  Skip (default)",
		"✅ Yes, continue",
	}

	prompt := &survey.Select{
		Message: message,
		Options: options,
		Default: options[0], // Default to Skip
	}

	var result string
	if err := survey.AskOne(prompt, &result); err != nil {
		return false, err
	}

	// Check if user selected "Yes"
	return result == options[1], nil
}

// PromptSelect prompts for selecting one option
func PromptSelect(message string, options []string) (string, error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	var result string
	if err := survey.AskOne(prompt, &result); err != nil {
		return "", err
	}

	return result, nil
}

// PromptMultiSelect prompts for selecting multiple options
func PromptMultiSelect(message string, options []string) ([]string, error) {
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}

	var results []string
	if err := survey.AskOne(prompt, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// PRTypeOptions returns the standard PR type options
func PRTypeOptions() []string {
	return []string{
		"✨ feat: New feature",
		"🐛 fix: Bug fix",
		"📝 docs: Documentation",
		"💄 style: Formatting, missing semi colons, etc.",
		"♻️  refactor: Code refactoring",
		"⚡ perf: Performance improvements",
		"✅ test: Adding tests",
		"🔧 chore: Maintenance",
		"⏪ revert: Revert changes",
	}
}

// ExtractPRType extracts the PR type from the selected option
func ExtractPRType(option string) string {
	// 移除 emoji 和描述，只保留类型
	if len(option) < 3 {
		return option
	}

	// 找到第一个冒号
	for i, r := range option {
		if r == ':' {
			if i > 0 {
				return option[3:i] // 跳过 emoji
			}
		}
	}

	return option
}
