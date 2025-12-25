package prompts

import (
	"testing"
)

func TestGolangProvider(t *testing.T) {
	provider := NewGolangProvider()
	prompts := provider.GetPrompts()

	if len(prompts) != 1 {
		t.Fatalf("Expected 1 prompt, got %d", len(prompts))
	}

	prompt := prompts[0]
	if prompt.Name != "golang-context-rule" {
		t.Errorf("Expected prompt name 'golang-context-rule', got '%s'", prompt.Name)
	}

	if prompt.Description == "" {
		t.Error("Prompt description should not be empty")
	}

	if prompt.Content == "" {
		t.Error("Prompt content should not be empty")
	}

	if len(prompt.Arguments) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(prompt.Arguments))
	}
}
