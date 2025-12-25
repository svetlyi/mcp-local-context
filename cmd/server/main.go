package main

import (
	"log/slog"
	"os"

	"github.com/svetlyi/mcp-local-context/internal/config"
	"github.com/svetlyi/mcp-local-context/internal/prompts"
	"github.com/svetlyi/mcp-local-context/internal/rules"
	"github.com/svetlyi/mcp-local-context/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Warn("Failed to load config, using defaults", "error", err)
		cfg = config.DefaultConfig()
	}

	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())

	for _, promptDir := range cfg.CustomPromptDirs {
		ruleProviders, err := rules.LoadRulesFromDirectory(promptDir)
		if err != nil {
			slog.Warn("Failed to load prompts from directory", "directory", promptDir, "error", err)
			continue
		}
		for _, provider := range ruleProviders {
			registry.Register(provider)
		}
		slog.Info("Loaded custom prompts", "directory", promptDir, "count", len(ruleProviders))
	}

	srv := server.New(registry)
	if err := srv.Run(); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
