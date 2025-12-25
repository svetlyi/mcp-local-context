package logging

import (
	"io"
	"log/slog"
	"os"

	"github.com/svetlyi/mcp-local-context/internal/config"
)

func Setup(cfg *config.Config) (func(), error) {
	var logWriter io.Writer = os.Stderr
	var close func() = func() {}

	if cfg.LogFile != "" {
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		logWriter = logFile
		close = func() { logFile.Close() }
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
