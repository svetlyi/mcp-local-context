package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	customProviders, err := custom.LoadPromptsFromDirectories(cfg.CustomPromptDirs)
	if err != nil {
		slog.Error("Failed to load custom prompts", "error", err)
		os.Exit(1)
	}
	for _, provider := range customProviders {
		registry.Register(provider)
	}
	if len(customProviders) > 0 {
		slog.Info("Loaded custom prompts", "count", len(customProviders))
	}

	srv, err := server.New(registry)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutting down server...")
		cancel()
	}()

	if err := srv.Run(ctx); err != nil && err != context.Canceled {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}
