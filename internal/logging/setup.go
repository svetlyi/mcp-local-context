package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/svetlyi/mcp-local-context/internal/config"
)

func Setup(cfg *config.Config) (func(), error) {
	var logWriter io.Writer
	var close func() = func() {}

	if cfg.LogFile != "" {
		logDir := filepath.Dir(cfg.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logWriter = logFile
		close = func() { logFile.Close() }
	} else {
		// Create a temporary log file
		tmpFile, err := os.CreateTemp("", "mcp-local-context-*.log")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary log file: %w", err)
		}
		logWriter = tmpFile
		close = func() { tmpFile.Close() }
	}

	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(logWriter, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return close, nil
}
