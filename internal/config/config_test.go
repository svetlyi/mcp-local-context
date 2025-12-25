package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svetlyi/mcp-local-context/internal/config"
)

func setupTestHome(t *testing.T) (string, func()) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	cleanup := func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}

	return tmpDir, cleanup
}

func TestLoadConfig(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	configDir := filepath.Join(tmpDir, ".mcp-local-context")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	cfgPath := filepath.Join(configDir, "config.json")
	testConfig := &config.Config{
		LogLevel: "debug",
	}

	data, err := json.Marshal(testConfig)
	require.NoError(t, err)

	err = os.WriteFile(cfgPath, data, 0644)
	require.NoError(t, err)

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	assert.Equal(t, "info", cfg.LogLevel)
	assert.NotNil(t, cfg.CustomRuleDirs)
}

func TestLoadConfigWithCustomPromptDirs(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	configDir := filepath.Join(tmpDir, ".mcp-local-context")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	cfgPath := filepath.Join(configDir, "config.json")
	customDir1 := filepath.Join(tmpDir, "custom1")
	customDir2 := filepath.Join(tmpDir, "custom2")
	testConfig := &config.Config{
		LogLevel:       "debug",
		CustomRuleDirs: []string{customDir1, customDir2},
	}

	data, err := json.Marshal(testConfig)
	require.NoError(t, err)

	err = os.WriteFile(cfgPath, data, 0644)
	require.NoError(t, err)

	cfg, err := config.Load()
	require.NoError(t, err)

	expectedDefaultRulesDir := filepath.Join(configDir, "rules")
	assert.GreaterOrEqual(t, len(cfg.CustomRuleDirs), 1, "Expected at least 1 directory (default rules dir)")

	foundDefault := false
	foundCustom1 := false
	foundCustom2 := false

	for _, dir := range cfg.CustomRuleDirs {
		if dir == expectedDefaultRulesDir {
			foundDefault = true
		}
		if dir == customDir1 {
			foundCustom1 = true
		}
		if dir == customDir2 {
			foundCustom2 = true
		}
	}

	assert.True(t, foundDefault, "Expected default rules dir %s to be in CustomRuleDirs, got %v", expectedDefaultRulesDir, cfg.CustomRuleDirs)
	assert.True(t, foundCustom1, "Expected custom dir %s to be in CustomRuleDirs, got %v", customDir1, cfg.CustomRuleDirs)
	assert.True(t, foundCustom2, "Expected custom dir %s to be in CustomRuleDirs, got %v", customDir2, cfg.CustomRuleDirs)
}

func TestLoadConfigWithoutCustomPromptDirs(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	configDir := filepath.Join(tmpDir, ".mcp-local-context")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	cfgPath := filepath.Join(configDir, "config.json")
	testConfig := &config.Config{}

	data, err := json.Marshal(testConfig)
	require.NoError(t, err)

	err = os.WriteFile(cfgPath, data, 0644)
	require.NoError(t, err)

	cfg, err := config.Load()
	require.NoError(t, err)

	expectedRulesDir := filepath.Join(configDir, "rules")
	assert.Equal(t, 1, len(cfg.CustomRuleDirs), "Expected 1 default custom prompt directory")
	assert.Equal(t, expectedRulesDir, cfg.CustomRuleDirs[0])
}
