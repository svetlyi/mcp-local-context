package server

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/svetlyi/mcp-local-context/internal/prompts"
)

func TestServerHandleInitialize(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	result, err := server.handleInitialize(nil)
	if err != nil {
		t.Fatalf("handleInitialize failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	if resultMap["protocolVersion"] != "2024-11-05" {
		t.Errorf("Expected protocolVersion '2024-11-05', got '%v'", resultMap["protocolVersion"])
	}

	serverInfo, ok := resultMap["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("serverInfo should be a map")
	}

	if serverInfo["name"] != "mcp-local-context" {
		t.Errorf("Expected server name 'mcp-local-context', got '%v'", serverInfo["name"])
	}
}

func TestServerHandlePromptsList(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	result, err := server.handlePromptsList(nil)
	if err != nil {
		t.Fatalf("handlePromptsList failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	prompts, ok := resultMap["prompts"].([]map[string]interface{})
	if !ok {
		t.Fatal("prompts should be a slice")
	}

	if len(prompts) < 1 {
		t.Error("Expected at least 1 prompt")
	}

	foundGolang := false
	for _, prompt := range prompts {
		if prompt["name"] == "golang-context-rule" {
			foundGolang = true
			break
		}
	}

	if !foundGolang {
		t.Error("Expected to find 'golang-context-rule' prompt")
	}
}

func TestServerHandlePromptsGet(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	params := map[string]interface{}{
		"name": "golang-context-rule",
	}
	paramsJSON, _ := json.Marshal(params)

	result, err := server.handlePromptsGet(paramsJSON)
	if err != nil {
		t.Fatalf("handlePromptsGet failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	if resultMap["name"] != "golang-context-rule" {
		t.Errorf("Expected name 'golang-context-rule', got '%v'", resultMap["name"])
	}

	messages, ok := resultMap["messages"].([]map[string]interface{})
	if !ok {
		t.Fatal("messages should be a slice")
	}

	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if messages[0]["role"] != "user" {
		t.Errorf("Expected role 'user', got '%v'", messages[0]["role"])
	}
}

func TestServerHandlePromptsGetNotFound(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	params := map[string]interface{}{
		"name": "nonexistent-prompt",
	}
	paramsJSON, _ := json.Marshal(params)

	_, err := server.handlePromptsGet(paramsJSON)
	if err == nil {
		t.Fatal("Expected error for nonexistent prompt")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' in error, got: %v", err)
	}
}

func TestServerHandleRequest(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "prompts/list",
		Params:  nil,
	}

	resp := server.handleRequest(req)
	if resp.Error != nil {
		t.Fatalf("Expected no error, got: %v", resp.Error)
	}

	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %v", resp.ID)
	}
}

func TestServerHandleRequestInvalidMethod(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "invalid/method",
		Params:  nil,
	}

	resp := server.handleRequest(req)
	if resp.Error == nil {
		t.Fatal("Expected error for invalid method")
	}

	if resp.Error.Code != -32601 {
		t.Errorf("Expected error code -32601, got %d", resp.Error.Code)
	}
}

func TestServerHandleRequestWithParams(t *testing.T) {
	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())
	server := New(registry)

	params := map[string]interface{}{
		"name": "golang-context-rule",
	}
	paramsJSON, _ := json.Marshal(params)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "prompts/get",
		Params:  paramsJSON,
	}

	resp := server.handleRequest(req)
	if resp.Error != nil {
		t.Fatalf("Expected no error, got: %v", resp.Error)
	}

	if resp.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %v", resp.ID)
	}
}
