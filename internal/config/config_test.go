package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Address != "localhost" {
		t.Errorf("Expected address 'localhost', got '%s'", cfg.Address)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.LogLevel)
	}
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Failed to get config directory: %v", err)
	}
	if configDir == "" {
		t.Error("Config directory should not be empty")
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	os.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".mcp-local-context")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	cfgPath := filepath.Join(configDir, "config.json")
	testConfig := &Config{
		Address:  "127.0.0.1",
		Port:     9000,
		LogLevel: "debug",
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Address != "127.0.0.1" {
		t.Errorf("Expected address '127.0.0.1', got '%s'", cfg.Address)
	}
	if cfg.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.LogLevel)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Should not fail when config file is missing: %v", err)
	}

	if cfg.Address != "localhost" {
		t.Errorf("Expected default address 'localhost', got '%s'", cfg.Address)
	}
}

func TestGetRulesDir(t *testing.T) {
	rulesDir, err := GetRulesDir()
	if err != nil {
		t.Fatalf("Failed to get rules directory: %v", err)
	}
	if rulesDir == "" {
		t.Error("Rules directory should not be empty")
	}
}

func TestLoadConfigWithCustomPromptDirs(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	os.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".mcp-local-context")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	cfgPath := filepath.Join(configDir, "config.json")
	customDir1 := filepath.Join(tmpDir, "custom1")
	customDir2 := filepath.Join(tmpDir, "custom2")
	testConfig := &Config{
		Address:          "127.0.0.1",
		Port:             9000,
		LogLevel:         "debug",
		CustomPromptDirs: []string{customDir1, customDir2},
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.CustomPromptDirs) != 2 {
		t.Errorf("Expected 2 custom prompt directories, got %d", len(cfg.CustomPromptDirs))
	}

	if cfg.CustomPromptDirs[0] != customDir1 {
		t.Errorf("Expected first directory %s, got %s", customDir1, cfg.CustomPromptDirs[0])
	}

	if cfg.CustomPromptDirs[1] != customDir2 {
		t.Errorf("Expected second directory %s, got %s", customDir2, cfg.CustomPromptDirs[1])
	}
}

func TestLoadConfigWithoutCustomPromptDirs(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	os.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".mcp-local-context")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	cfgPath := filepath.Join(configDir, "config.json")
	testConfig := &Config{
		Address: "127.0.0.1",
		Port:    9000,
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.CustomPromptDirs) != 1 {
		t.Errorf("Expected 1 default custom prompt directory, got %d", len(cfg.CustomPromptDirs))
	}

	expectedRulesDir, _ := GetRulesDir()
	if cfg.CustomPromptDirs[0] != expectedRulesDir {
		t.Errorf("Expected default rules directory %s, got %s", expectedRulesDir, cfg.CustomPromptDirs[0])
	}
}
