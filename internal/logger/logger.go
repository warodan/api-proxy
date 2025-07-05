package logger

import (
	"api-proxy/internal/config"
	"log/slog"
	"os"
)

func New(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch cfg.LoggerLevel {
	case "DEBUG":
		opts.Level = slog.LevelDebug
	case "WARN":
		opts.Level = slog.LevelWarn
	case "ERROR":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
