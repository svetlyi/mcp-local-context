package custom

import (
	"bytes"
	"fmt"
	"log/slog"
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
			return nil, fmt.Errorf("failed to read prompt file %s: %w", entry.Name(), err)
		}

		if len(content) == 0 {
			slog.Warn("Empty prompt file, skipping", "file", entry.Name())
			continue
		}

		config, cleanContent := parseConfig(content)
		title := config["title"]
		language := config["lang"]

		description := title
		if description == "" {
			description = extractDescription(cleanContent)
		}
		if description == "" {
			description = fmt.Sprintf("Custom prompt: %s", strings.TrimSuffix(entry.Name(), ".md"))
		}

		if len(cleanContent) == 0 {
			slog.Warn("Prompt file has no content after parsing config, skipping", "file", entry.Name())
			continue
		}

		promptName := strings.TrimSuffix(entry.Name(), ".md")
		if promptName == "" {
			slog.Warn("Invalid prompt filename, skipping", "file", entry.Name())
			continue
		}

		provider := newCustomPromptProvider(promptName, string(cleanContent), description, language)
		providers = append(providers, provider)
	}

	return providers, nil
}

// parseConfig parses key:value configuration from the beginning of content.
// Returns a map of config values and the remaining content without the config section.
// Config lines must be in the format "key:value" and appear at the start of the file.
// The config section ends when a blank line or non-config line is encountered.
func parseConfig(content []byte) (map[string]string, []byte) {
	config := make(map[string]string)
	lines := bytes.Split(content, []byte("\n"))

	var configEndIdx int
	for i, line := range lines {
		lineStr := strings.TrimSpace(string(line))

		if lineStr == "" {
			configEndIdx = i
			break
		}

		parts := strings.SplitN(lineStr, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				config[key] = value
				configEndIdx = i + 1
				continue
			}
		}

		configEndIdx = i
		break
	}

	if configEndIdx > 0 {
		remainingLines := lines[configEndIdx:]
		cleanContent := bytes.Join(remainingLines, []byte("\n"))
		if len(cleanContent) > 0 && cleanContent[0] == '\n' {
			cleanContent = cleanContent[1:]
		}
		return config, cleanContent
	}

	return config, content
}

// extractDescription extracts the description from the first line of content.
// It removes markdown heading markers (#, ##, etc.) if present.
// Handles both Unix (\n) and Windows (\r\n) line endings.
func extractDescription(content []byte) string {
	// Find the first newline or use the entire content if no newline exists
	newlineIdx := bytes.IndexByte(content, '\n')
	var firstLine []byte
	if newlineIdx >= 0 {
		firstLine = content[:newlineIdx]
		// Remove Windows carriage return if present
		if len(firstLine) > 0 && firstLine[len(firstLine)-1] == '\r' {
			firstLine = firstLine[:len(firstLine)-1]
		}
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

func LoadPromptsFromDirectories(dirs []string) ([]prompts.Provider, error) {
	var allProviders []prompts.Provider

	for _, dir := range dirs {
		providers, err := LoadPromptsFromDirectory(dir)
		if err != nil {
			slog.Warn("Failed to load prompts from directory, skipping", "directory", dir, "error", err)
			continue
		}
		allProviders = append(allProviders, providers...)
	}

	return allProviders, nil
}

type customPromptProvider struct {
	name        string
	content     string
	description string
	language    string
}

func newCustomPromptProvider(name, content, description, language string) *customPromptProvider {
	return &customPromptProvider{
		name:        name,
		content:     content,
		description: description,
		language:    language,
	}
}

func (c *customPromptProvider) GetPrompts() []prompts.Prompt {
	return []prompts.Prompt{
		{
			Name:        c.name,
			Description: c.description,
			Arguments:   []prompts.PromptArgument{},
			Content:     c.content,
			Language:    c.language,
		},
	}
}
