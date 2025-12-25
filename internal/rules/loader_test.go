package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

func TestLoadRulesFromDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	rule1Path := filepath.Join(tmpDir, "rule1.md")
	rule1Content := "# Rule 1\n\nThis is rule 1 content."
	if err := os.WriteFile(rule1Path, []byte(rule1Content), 0644); err != nil {
		t.Fatalf("Failed to write rule1: %v", err)
	}

	rule2Path := filepath.Join(tmpDir, "rule2.md")
	rule2Content := "# Rule 2\n\nThis is rule 2 content."
	if err := os.WriteFile(rule2Path, []byte(rule2Content), 0644); err != nil {
		t.Fatalf("Failed to write rule2: %v", err)
	}

	nonMdPath := filepath.Join(tmpDir, "not-a-rule.txt")
	if err := os.WriteFile(nonMdPath, []byte("not a markdown file"), 0644); err != nil {
		t.Fatalf("Failed to write non-md file: %v", err)
	}

	providers, err := LoadRulesFromDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	if len(providers) != 2 {
		t.Fatalf("Expected 2 providers, got %d", len(providers))
	}

	allPrompts := make([]prompts.Prompt, 0)
	for _, provider := range providers {
		allPrompts = append(allPrompts, provider.GetPrompts()...)
	}

	if len(allPrompts) != 2 {
		t.Fatalf("Expected 2 prompts, got %d", len(allPrompts))
	}

	foundRule1 := false
	foundRule2 := false
	for _, prompt := range allPrompts {
		if prompt.Name == "rule1" {
			foundRule1 = true
			if prompt.Content != rule1Content {
				t.Errorf("Rule1 content mismatch")
			}
		}
		if prompt.Name == "rule2" {
			foundRule2 = true
			if prompt.Content != rule2Content {
				t.Errorf("Rule2 content mismatch")
			}
		}
	}

	if !foundRule1 {
		t.Error("Rule1 not found")
	}
	if !foundRule2 {
		t.Error("Rule2 not found")
	}
}

func TestLoadRulesFromDirectoryMissing(t *testing.T) {
	tmpDir := t.TempDir()
	nonexistentDir := filepath.Join(tmpDir, "nonexistent")

	providers, err := LoadRulesFromDirectory(nonexistentDir)
	if err != nil {
		t.Fatalf("Should not fail for nonexistent directory: %v", err)
	}

	if len(providers) != 0 {
		t.Errorf("Expected 0 providers, got %d", len(providers))
	}
}

func TestLoadRulesFromDirectoryEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	providers, err := LoadRulesFromDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load rules from empty directory: %v", err)
	}

	if len(providers) != 0 {
		t.Errorf("Expected 0 providers, got %d", len(providers))
	}
}

func TestRuleProvider(t *testing.T) {
	content := "# Test Rule\n\nThis is test content."
	provider := NewRuleProvider("test-rule", content)

	prompts := provider.GetPrompts()
	if len(prompts) != 1 {
		t.Fatalf("Expected 1 prompt, got %d", len(prompts))
	}

	prompt := prompts[0]
	if prompt.Name != "test-rule" {
		t.Errorf("Expected name 'test-rule', got '%s'", prompt.Name)
	}

	if prompt.Content != content {
		t.Errorf("Content mismatch")
	}
}
