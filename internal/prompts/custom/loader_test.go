package custom_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svetlyi/mcp-local-context/internal/prompts"
	"github.com/svetlyi/mcp-local-context/internal/prompts/custom"
)

func TestLoadPromptsFromDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	prompt1Path := filepath.Join(tmpDir, "prompt1.md")
	prompt1Content := `lang:golang
title:Prompt 1 Title

# Prompt 1
This is prompt 1 content.`
	err := os.WriteFile(prompt1Path, []byte(prompt1Content), 0644)
	require.NoError(t, err)

	prompt2Path := filepath.Join(tmpDir, "prompt2.md")
	prompt2Content := `# Prompt 2 Header
This is prompt 2 content.`
	err = os.WriteFile(prompt2Path, []byte(prompt2Content), 0644)
	require.NoError(t, err)

	nonMdPath := filepath.Join(tmpDir, "not-a-prompt.txt")
	err = os.WriteFile(nonMdPath, []byte("not a markdown file"), 0644)
	require.NoError(t, err)

	providers, err := custom.LoadPromptsFromDirectory(tmpDir)
	require.NoError(t, err)
	assert.Len(t, providers, 2, "Expected 2 providers")

	allPrompts := make([]prompts.Prompt, 0)
	for _, provider := range providers {
		allPrompts = append(allPrompts, provider.GetPrompts()...)
	}

	assert.Len(t, allPrompts, 2, "Expected 2 prompts")

	foundPrompt1 := false
	foundPrompt2 := false
	for _, prompt := range allPrompts {
		if prompt.Name == "prompt1" {
			foundPrompt1 = true
			assert.Contains(t, prompt.Content, "# Prompt 1", "Prompt1 content should contain the actual content")
			assert.Contains(t, prompt.Content, "This is prompt 1 content", "Prompt1 content should contain the actual content")
			assert.NotContains(t, prompt.Content, "lang:golang", "Prompt1 content should not contain config")
			assert.NotContains(t, prompt.Content, "title:", "Prompt1 content should not contain config")
			assert.Equal(t, "Prompt 1 Title", prompt.Description, "Prompt1 description should be extracted from title key")
			assert.Equal(t, "golang", prompt.Language, "Prompt1 language should be extracted from lang key")
		}
		if prompt.Name == "prompt2" {
			foundPrompt2 = true
			assert.Equal(t, prompt2Content, prompt.Content, "Prompt2 content mismatch")
			assert.Equal(t, "Prompt 2 Header", prompt.Description, "Prompt2 description should be extracted from first line")
			assert.Equal(t, "", prompt.Language, "Prompt2 should have no language")
		}
	}

	assert.True(t, foundPrompt1, "Prompt1 not found")
	assert.True(t, foundPrompt2, "Prompt2 not found")
}

func TestLoadPromptsFromDirectoryMissing(t *testing.T) {
	tmpDir := t.TempDir()
	nonexistentDir := filepath.Join(tmpDir, "nonexistent")

	providers, err := custom.LoadPromptsFromDirectory(nonexistentDir)
	require.NoError(t, err, "Should not fail for nonexistent directory")
	assert.Len(t, providers, 0, "Expected 0 providers")
}

func TestLoadPromptsFromDirectoryEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	providers, err := custom.LoadPromptsFromDirectory(tmpDir)
	require.NoError(t, err, "Failed to load prompts from empty directory")
	assert.Len(t, providers, 0, "Expected 0 providers")
}
