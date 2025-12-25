package server

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

type Server struct {
	registry  *prompts.Registry
	mcpServer *mcp.Server
}

func New(registry *prompts.Registry) *Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-local-context",
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

	return s
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

// Run starts the MCP server on stdio transport
func (s *Server) Run() error {
	ctx := context.Background()
	transport := &mcp.StdioTransport{}
	return s.mcpServer.Run(ctx, transport)
}
