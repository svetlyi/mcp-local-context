package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultConfigDir  = ".mcp-local-context"
	defaultRulesDir   = "rules"
	defaultConfigFile = "config.json"
)

type Config struct {
	LogLevel       string   `json:"log_level,omitempty"`
	LogFile        string   `json:"log_file,omitempty"`
	CustomRuleDirs []string `json:"custom_rule_dirs,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		LogLevel:       "info",
		CustomRuleDirs: make([]string, 0),
	}
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, defaultConfigDir), nil
}

func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, defaultConfigFile), nil
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()

	defer func() {
		defaultRulesDir, err := getRulesDir()
		if err == nil {
			config.CustomRuleDirs = append(config.CustomRuleDirs, defaultRulesDir)
		}
	}()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	for i, dir := range config.CustomRuleDirs {
		expanded := expandPath(dir)
		config.CustomRuleDirs[i] = expanded
	}

	config.LogFile = expandPath(config.LogFile)

	return config, nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}

func getRulesDir() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, defaultRulesDir), nil
}
