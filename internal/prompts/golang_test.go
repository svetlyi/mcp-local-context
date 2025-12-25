package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolangProvider(t *testing.T) {
	provider := NewGolangProvider()
	prompts := provider.GetPrompts()

	require.Len(t, prompts, 1, "Expected 1 prompt")

	prompt := prompts[0]
	assert.Equal(t, "golang-context-rule", prompt.Name)
	assert.NotEmpty(t, prompt.Description, "Prompt description should not be empty")
	assert.NotEmpty(t, prompt.Content, "Prompt content should not be empty")
	assert.Len(t, prompt.Arguments, 0, "Expected 0 arguments")
}
