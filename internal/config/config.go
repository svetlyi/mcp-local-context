package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Address          string   `json:"address,omitempty"`
	Port             int      `json:"port,omitempty"`
	LogLevel         string   `json:"log_level,omitempty"`
	CustomPromptDirs []string `json:"custom_prompt_dirs,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Address:  "localhost",
		Port:     8080,
		LogLevel: "info",
	}
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".mcp-local-context"), nil
}

func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()

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

	if len(config.CustomPromptDirs) == 0 {
		rulesDir, err := GetRulesDir()
		if err == nil {
			config.CustomPromptDirs = []string{rulesDir}
		}
	} else {
		expandedDirs := make([]string, 0, len(config.CustomPromptDirs))
		for _, dir := range config.CustomPromptDirs {
			expanded := expandPath(dir)
			expandedDirs = append(expandedDirs, expanded)
		}
		config.CustomPromptDirs = expandedDirs
	}

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

func GetRulesDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "rules"), nil
}
