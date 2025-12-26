package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

type Server struct {
	registry  *prompts.Registry
	mcpServer *mcp.Server
}

func New(registry *prompts.Registry) (*Server, error) {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-local-context",
		Title:   "Local Context Instructions Server",
		Version: "0.1.0",
	}, nil)

	s := &Server{
		registry:  registry,
		mcpServer: mcpServer,
	}

	allPrompts := registry.GetAllPrompts()
	for _, prompt := range allPrompts {
		s.registerPrompt(prompt)
	}

	if err := s.registerTools(); err != nil {
		slog.Error("Failed to register tools", "error", err)
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return s, nil
}

func (s *Server) registerPrompt(prompt prompts.Prompt) {
	// Convert our Prompt to MCP Prompt
	mcpPrompt := &mcp.Prompt{
		Name:        prompt.Name,
		Description: prompt.Description,
	}

	// Convert arguments
	if len(prompt.Arguments) > 0 {
		mcpPrompt.Arguments = make([]*mcp.PromptArgument, 0, len(prompt.Arguments))
		for _, arg := range prompt.Arguments {
			mcpPrompt.Arguments = append(mcpPrompt.Arguments, &mcp.PromptArgument{
				Name:        arg.Name,
				Description: arg.Description,
				Required:    arg.Required,
			})
		}
	}

	// Create handler that returns the prompt content
	handler := func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		p := s.registry.GetPrompt(req.Params.Name)
		if p == nil {
			return nil, fmt.Errorf("prompt not found: %s", req.Params.Name)
		}

		// Convert to MCP message format
		messages := []*mcp.PromptMessage{
			{
				Role:    mcp.Role("user"),
				Content: &mcp.TextContent{Text: p.Content},
			},
		}

		return &mcp.GetPromptResult{
			Messages: messages,
		}, nil
	}

	s.mcpServer.AddPrompt(mcpPrompt, handler)
}

func (s *Server) registerTools() error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic during tool registration", "panic", r)
		}
	}()

	type getSupportedLanguagesArgs struct{}
	type getSupportedLanguagesOutput struct {
		Languages []string `json:"languages"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "list_supported_languages",
		Description: "CRITICAL: You MUST call this tool FIRST whenever you encounter third-party packages, libraries, dependencies, or external modules in ANY codebase. This tool reveals which programming languages have proven, battle-tested context instructions available. These instructions are ESSENTIAL for working correctly with dependencies - they prevent you from making incorrect assumptions, using outdated documentation, or misunderstanding package structures. Workflow: 1) Call this tool first to see available languages, 2) Then immediately call get_context_instructions with the relevant language. DO NOT attempt to work with third-party code without first checking and applying these context instructions - doing so leads to errors, incorrect implementations, and wasted time.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getSupportedLanguagesArgs) (*mcp.CallToolResult, getSupportedLanguagesOutput, error) {
		languages := s.registry.GetSupportedLanguages()
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Supported languages: %s", strings.Join(languages, ", "))},
			},
		}, getSupportedLanguagesOutput{Languages: languages}, nil
	})

	type getContextInstructionsArgs struct {
		// Language is the programming language of the third-party package/library you're working with. REQUIRED. Examples: "golang" for Go modules/packages, "python" for pip packages, "javascript" for npm packages, etc. Use list_supported_languages first to see available options.
		Language string `json:"language" jsonschema:"required"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_context_instructions",
		Description: "MANDATORY when working with ANY third-party package, library, dependency, or external module. This tool provides authoritative, systematic instructions for correctly understanding and working with code in a specific language ecosystem. These instructions are CRITICAL because: 1) They show you how to locate and inspect the ACTUAL source code in local caches (not outdated online docs), 2) They provide proven methods for understanding package structures, module systems, and dependency resolution, 3) They prevent common mistakes and incorrect assumptions that lead to bugs. BEFORE calling: Use list_supported_languages first to verify support. WHEN to use: IMMEDIATELY when you see imports, require statements, package declarations, or any reference to external code. These instructions are your source of truth - ignore them at your peril. They are specifically designed to help you work with real, local codebases rather than relying on potentially outdated or incomplete documentation.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getContextInstructionsArgs) (*mcp.CallToolResult, any, error) {
		if args.Language == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: language argument is required"},
				},
			}, nil, fmt.Errorf("language argument is required")
		}

		prompts := s.registry.GetPromptsByLanguage(args.Language)
		if len(prompts) == 0 {
			supportedLanguages := s.registry.GetSupportedLanguages()
			message := fmt.Sprintf("⚠️ No context instructions found for language: %s", args.Language)
			if len(supportedLanguages) > 0 {
				message += fmt.Sprintf("\n\n✅ Available languages with context instructions: %s\n\n", strings.Join(supportedLanguages, ", "))
				message += "⚠️ CRITICAL: You MUST use list_supported_languages FIRST before requesting instructions. "
				message += "Working with third-party packages without proper context instructions leads to errors and incorrect implementations. "
				message += "If your language is listed above, use it. If not, you may need to work without specialized instructions, but this is NOT recommended for third-party code."
			} else {
				message += "\n\n⚠️ No languages are currently supported. This means you cannot use specialized context instructions for third-party packages."
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: message},
				},
			}, nil, nil
		}

		var contents []mcp.Content
		for _, prompt := range prompts {
			contents = append(contents, &mcp.TextContent{
				Text: fmt.Sprintf("# %s\n\n%s", prompt.Name, prompt.Content),
			})
		}

		return &mcp.CallToolResult{
			Content: contents,
		}, nil, nil
	})

	return nil
}

func (s *Server) Run(ctx context.Context) error {
	transport := &mcp.StdioTransport{}
	return s.mcpServer.Run(ctx, transport)
}
