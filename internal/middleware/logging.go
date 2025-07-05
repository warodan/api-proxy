package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type ctxKeyLogger struct{}

func InjectLoggerMiddleware(baseLogger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			reqID := uuid.New().String()
			start := time.Now()

			reqLogger := baseLogger.With(
				slog.String("request_id", reqID),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
			)
			ctx := context.WithValue(req.Context(), ctxKeyLogger{}, reqLogger)

			next.ServeHTTP(writer, req.WithContext(ctx))

			duration := time.Since(start)
			reqLogger.Info("request completed", slog.Duration("duration", duration))
		})
	}
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxKeyLogger{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
