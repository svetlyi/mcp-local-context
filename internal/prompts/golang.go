package prompts

import _ "embed"

//go:embed golang.md
var golangPromptContent string

type GolangProvider struct{}

func NewGolangProvider() *GolangProvider {
	return &GolangProvider{}
}

func (g *GolangProvider) GetPrompts() []Prompt {
	return []Prompt{
		{
			Name:        "golang-context-rule",
			Description: "Provides a systematic approach for working with third-party Go packages by referencing the Go module cache",
			Arguments:   []PromptArgument{},
			Content:     golangPromptContent,
		},
	}
}
