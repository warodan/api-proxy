package main

import (
	"api-proxy/internal/config"
	"api-proxy/internal/handler"
	"api-proxy/internal/logger"
	"api-proxy/internal/middleware"
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	tempLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.Load(tempLogger)
	if err := cfg.Validate(); err != nil {
		tempLogger.Error("invalid config", "err", err)
		os.Exit(1)
	}

	logger := logger.New(cfg)
	client := resty.New()

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log := resp.Request.Context().Value("logger").(*slog.Logger)
		log.Info("Resty request completed", "status", resp.StatusCode(), "url", resp.Request.URL)
		return nil
	})

	mux := http.NewServeMux()
	postHandler := handler.NewPostHandler(client)
	mux.HandleFunc("/posts/", postHandler.ProxyPost)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.InjectLoggerMiddleware(logger)(mux),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	go func() {
		logger.Info("Server is starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	logger.Info("Graceful shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error", "err", err)
	} else {
		logger.Info("Server gracefully stopped")
	}
}
