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
	defaultPromptsDir = "prompts"
	defaultConfigFile = "config.json"
)

type Config struct {
	LogLevel         string   `json:"log_level,omitempty"`
	LogFile          string   `json:"log_file,omitempty"`
	CustomPromptDirs []string `json:"custom_prompt_dirs,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		LogLevel:         "info",
		CustomPromptDirs: make([]string, 0),
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
		defaultPromptsDir, err := getPromptsDir()
		if err == nil {
			config.CustomPromptDirs = append(config.CustomPromptDirs, defaultPromptsDir)
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

	for i, dir := range config.CustomPromptDirs {
		expanded := expandPath(dir)
		config.CustomPromptDirs[i] = expanded
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

func getPromptsDir() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, defaultPromptsDir), nil
}
