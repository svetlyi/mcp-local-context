package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	registry *prompts.Registry
}

func New(registry *prompts.Registry) *Server {
	return &Server{
		registry: registry,
	}
}

func (s *Server) handleInitialize(params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"prompts": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "mcp-local-context",
			"version": "0.1.0",
		},
	}, nil
}

func (s *Server) handlePromptsList(params json.RawMessage) (interface{}, error) {
	allPrompts := s.registry.GetAllPrompts()
	promptsList := make([]map[string]interface{}, 0, len(allPrompts))

	for _, prompt := range allPrompts {
		promptsList = append(promptsList, map[string]interface{}{
			"name":        prompt.Name,
			"description": prompt.Description,
			"arguments":   prompt.Arguments,
		})
	}

	return map[string]interface{}{
		"prompts": promptsList,
	}, nil
}

func (s *Server) handlePromptsGet(params json.RawMessage) (interface{}, error) {
	var requestParams struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(params, &requestParams); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	prompt := s.registry.GetPrompt(requestParams.Name)
	if prompt == nil {
		return nil, fmt.Errorf("prompt not found: %s", requestParams.Name)
	}

	return map[string]interface{}{
		"name":        prompt.Name,
		"description": prompt.Description,
		"arguments":   prompt.Arguments,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt.Content,
			},
		},
	}, nil
}

func (s *Server) handleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	var result interface{}
	var err error

	switch req.Method {
	case "initialize":
		result, err = s.handleInitialize(req.Params)
	case "prompts/list":
		result, err = s.handlePromptsList(req.Params)
	case "prompts/get":
		result, err = s.handlePromptsGet(req.Params)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}

	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *Server) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			slog.Warn("Failed to parse request", "error", err)
			continue
		}

		resp := s.handleRequest(&req)
		
		if req.ID == nil {
			continue
		}

		if err := encoder.Encode(resp); err != nil {
			slog.Warn("Failed to encode response", "error", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	return nil
}
