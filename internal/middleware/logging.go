package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func InjectLoggerMiddleware(baseLogger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			reqId := uuid.New().String()
			reqLogger := baseLogger.With(
				slog.String("request_id", reqId),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
			)
			ctx := context.WithValue(req.Context(), "logger", reqLogger)

			next.ServeHTTP(writer, req.WithContext(ctx))
		})
	}
}
