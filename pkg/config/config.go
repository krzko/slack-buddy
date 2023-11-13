package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	configDir  = ".slack-buddy"
	configFile = "config.yaml"
)

type Config struct {
	ApiToken    string       `yaml:"api_token"`
	UserId      string       `yaml:"user_id"`
	DisplayName string       `yaml:"display_name"`
	CustomItems []CustomItem `yaml:"custom_items"`
}

type CustomItem struct {
	Title       string   `yaml:"title"`
	Tooltip     string   `yaml:"tooltip"`
	StatusText  string   `yaml:"status_text"`
	StatusEmoji string   `yaml:"status_emoji"`
	Days        []string `yaml:"days"`
	StartTime   string   `yaml:"start_time"`
	EndTime     string   `yaml:"end_time"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(homeDir, configDir, configFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.ApiToken == "" {
		config.ApiToken = "your_default_api_token"
	}
	if config.UserId == "" {
		config.UserId = "your_default_user_id"
	}
	if config.DisplayName == "" {
		config.DisplayName = "your_default_display_name"
	}

	return &config, nil
}

func (c *Config) SaveConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, configDir, configFile)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, data, 0644)
	return err
}
