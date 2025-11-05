package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Email              string `mapstructure:"email"`
	JiraAPIToken       string `mapstructure:"jira_api_token"`
	JiraServiceAddress string `mapstructure:"jira_service_address"`
	GitHubToken        string `mapstructure:"github_token"`
	BranchPrefix       string `mapstructure:"branch_prefix"`
	OpenAIKey          string `mapstructure:"openai_key"`
	DeepSeekKey        string `mapstructure:"deepseek_key"`
	OpenAIProxyURL     string `mapstructure:"openai_proxy_url"`
	OpenAIProxyKey     string `mapstructure:"openai_proxy_key"`
}

var globalConfig *Config

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// 设置配置文件路径
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "quick-workflow")
	configFile := filepath.Join(configDir, "config.yaml")

	// 创建配置目录（如果不存在）
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// 环境变量前缀
	viper.SetEnvPrefix("QK")
	viper.AutomaticEnv()

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在是正常的，将创建默认配置
	}

	// 设置默认值
	setDefaults()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	globalConfig = &cfg
	return globalConfig, nil
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		// 如果还没加载，尝试加载
		cfg, _ := Load()
		return cfg
	}
	return globalConfig
}

// Save saves the current configuration to file
func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "quick-workflow")
	configFile := filepath.Join(configDir, "config.yaml")

	// 创建配置目录
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 设置值
	viper.Set("email", cfg.Email)
	viper.Set("jira_api_token", cfg.JiraAPIToken)
	viper.Set("jira_service_address", cfg.JiraServiceAddress)
	viper.Set("github_token", cfg.GitHubToken)
	viper.Set("branch_prefix", cfg.BranchPrefix)
	viper.Set("openai_key", cfg.OpenAIKey)
	viper.Set("deepseek_key", cfg.DeepSeekKey)
	viper.Set("openai_proxy_url", cfg.OpenAIProxyURL)
	viper.Set("openai_proxy_key", cfg.OpenAIProxyKey)

	// 写入文件
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	globalConfig = cfg
	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Email == "" {
		return fmt.Errorf("email is required")
	}
	if c.JiraAPIToken == "" {
		return fmt.Errorf("jira_api_token is required")
	}
	if c.JiraServiceAddress == "" {
		return fmt.Errorf("jira_service_address is required")
	}
	if c.GitHubToken == "" {
		return fmt.Errorf("github_token is required (run: gh auth token)")
	}
	return nil
}

// IsConfigured checks if the application is configured
func IsConfigured() bool {
	cfg := Get()
	if cfg == nil {
		return false
	}
	return cfg.Validate() == nil
}

func setDefaults() {
	viper.SetDefault("branch_prefix", "")
}

