package logger

import (
	"api-proxy/internal/middleware"
	"context"
	"log/slog"
)

func MustLoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(middleware.CtxKeyLogger{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
