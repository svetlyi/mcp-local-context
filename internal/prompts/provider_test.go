package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	provider1 := NewGolangProvider()
	registry.Register(provider1)

	prompts := registry.GetAllPrompts()
	require.Len(t, prompts, 1, "Expected 1 prompt")

	prompt := registry.GetPrompt("golang-context-rule")
	require.NotNil(t, prompt, "Expected to find 'golang-context-rule' prompt")
	assert.Equal(t, "golang-context-rule", prompt.Name)

	nonexistent := registry.GetPrompt("nonexistent")
	assert.Nil(t, nonexistent, "Expected nil for nonexistent prompt")
}

func TestRegistryMultipleProviders(t *testing.T) {
	registry := NewRegistry()

	provider1 := NewGolangProvider()
	registry.Register(provider1)
	registry.Register(provider1)

	prompts := registry.GetAllPrompts()
	assert.Len(t, prompts, 2, "Expected 2 prompts")
}
