package main

import (
	"log/slog"
	"os"

	"github.com/svetlyi/mcp-local-context/internal/config"
	"github.com/svetlyi/mcp-local-context/internal/logging"
	"github.com/svetlyi/mcp-local-context/internal/prompts"
	"github.com/svetlyi/mcp-local-context/internal/prompts/custom"
	"github.com/svetlyi/mcp-local-context/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Warn("Failed to load config, using defaults", "error", err)
		cfg = config.DefaultConfig()
	}

	closeLog, err := logging.Setup(cfg)
	if err != nil {
		slog.Error("Failed to setup logging", "error", err)
		os.Exit(1)
	}
	defer closeLog()

	registry := prompts.NewRegistry()
	registry.Register(prompts.NewGolangProvider())

	customProviders := custom.LoadPromptsFromDirectories(cfg.CustomPromptDirs)
	for _, provider := range customProviders {
		registry.Register(provider)
	}
	if len(customProviders) > 0 {
		slog.Info("Loaded custom prompts", "count", len(customProviders))
	}

	srv := server.New(registry)
	if err := srv.Run(); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
