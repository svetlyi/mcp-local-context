package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

func LoadRulesFromDirectory(rulesDir string) ([]prompts.Provider, error) {
	var providers []prompts.Provider

	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return providers, nil
	}

	entries, err := os.ReadDir(rulesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}

		rulePath := filepath.Join(rulesDir, entry.Name())
		content, err := os.ReadFile(rulePath)
		if err != nil {
			continue
		}

		ruleName := strings.TrimSuffix(entry.Name(), ".md")
		provider := NewRuleProvider(ruleName, string(content))
		providers = append(providers, provider)
	}

	return providers, nil
}

type RuleProvider struct {
	name    string
	content string
}

func NewRuleProvider(name, content string) *RuleProvider {
	return &RuleProvider{
		name:    name,
		content: content,
	}
}

func (r *RuleProvider) GetPrompts() []prompts.Prompt {
	return []prompts.Prompt{
		{
			Name:        r.name,
			Description: fmt.Sprintf("Custom rule loaded from rules directory: %s", r.name),
			Arguments:   []prompts.PromptArgument{},
			Content:     r.content,
		},
	}
}
