package prompts

import (
	"testing"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	provider1 := NewGolangProvider()
	registry.Register(provider1)

	prompts := registry.GetAllPrompts()
	if len(prompts) != 1 {
		t.Fatalf("Expected 1 prompt, got %d", len(prompts))
	}

	prompt := registry.GetPrompt("golang-context-rule")
	if prompt == nil {
		t.Fatal("Expected to find 'golang-context-rule' prompt")
	}

	if prompt.Name != "golang-context-rule" {
		t.Errorf("Expected prompt name 'golang-context-rule', got '%s'", prompt.Name)
	}

	nonexistent := registry.GetPrompt("nonexistent")
	if nonexistent != nil {
		t.Error("Expected nil for nonexistent prompt")
	}
}

func TestRegistryMultipleProviders(t *testing.T) {
	registry := NewRegistry()

	provider1 := NewGolangProvider()
	registry.Register(provider1)
	registry.Register(provider1)

	prompts := registry.GetAllPrompts()
	if len(prompts) != 2 {
		t.Fatalf("Expected 2 prompts, got %d", len(prompts))
	}
}
