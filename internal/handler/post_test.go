package handler_test

import (
	"api-proxy/internal/handler"
	"api-proxy/internal/middleware"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProxyPost(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client := resty.New()

	h := handler.NewPostHandler(client)

	mux := http.NewServeMux()
	mux.HandleFunc("/posts/", h.ProxyPost)

	wrapped := middleware.InjectLoggerMiddleware(logger)(mux)

	testServer := httptest.NewServer(wrapped)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/posts/1")
	if err != nil {
		t.Fatalf("Request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("We were expecting 200 status, we got %d", resp.StatusCode)
	}
}
