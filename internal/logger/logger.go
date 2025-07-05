package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	opts := &slog.HandlerOptions{}

	opts.Level = slog.LevelInfo

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
