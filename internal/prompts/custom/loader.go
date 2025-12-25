package custom

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

func LoadPromptsFromDirectory(promptsDir string) ([]prompts.Provider, error) {
	var providers []prompts.Provider

	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		return providers, nil
	}

	entries, err := os.ReadDir(promptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}

		promptPath := filepath.Join(promptsDir, entry.Name())
		content, err := os.ReadFile(promptPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read prompt file %s: %w", promptPath, err)
		}

		description := extractDescription(content)
		if description == "" {
			description = fmt.Sprintf("Custom prompt: %s", strings.TrimSuffix(entry.Name(), ".md"))
		}

		promptName := strings.TrimSuffix(entry.Name(), ".md")
		provider := newCustomPromptProvider(promptName, string(content), description)
		providers = append(providers, provider)
	}

	return providers, nil
}

// extractDescription extracts the description from the first line of content.
// It removes markdown heading markers (#, ##, etc.) if present.
func extractDescription(content []byte) string {
	// Find the first newline or use the entire content if no newline exists
	newlineIdx := bytes.IndexByte(content, '\n')
	var firstLine []byte
	if newlineIdx >= 0 {
		firstLine = content[:newlineIdx]
	} else {
		firstLine = content
	}

	// Trim whitespace and remove markdown heading markers
	firstLineStr := strings.TrimSpace(string(firstLine))
	if firstLineStr == "" {
		return ""
	}

	// Remove markdown heading markers (#, ##, etc.)
	description := strings.TrimLeft(firstLineStr, "# ")
	if description == "" {
		description = firstLineStr
	}

	return description
}

func LoadPromptsFromDirectories(dirs []string) []prompts.Provider {
	var allProviders []prompts.Provider

	for _, dir := range dirs {
		providers, err := LoadPromptsFromDirectory(dir)
		if err != nil {
			continue
		}
		allProviders = append(allProviders, providers...)
	}

	return allProviders
}

type customPromptProvider struct {
	name        string
	content     string
	description string
}

func newCustomPromptProvider(name, content, description string) *customPromptProvider {
	return &customPromptProvider{
		name:        name,
		content:     content,
		description: description,
	}
}

func (c *customPromptProvider) GetPrompts() []prompts.Prompt {
	return []prompts.Prompt{
		{
			Name:        c.name,
			Description: c.description,
			Arguments:   []prompts.PromptArgument{},
			Content:     c.content,
		},
	}
}
