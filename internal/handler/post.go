package handler

import (
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"strings"
)

type PostHandler struct {
	client *resty.Client
}

func NewPostHandler(client *resty.Client) *PostHandler {
	return &PostHandler{client: client}
}

func (handler *PostHandler) ProxyPost(writer http.ResponseWriter, req *http.Request) {
	logger := req.Context().Value("logger").(*slog.Logger)

	id := strings.TrimPrefix(req.URL.Path, "/posts/")

	resp, err := handler.client.R().
		SetHeader("Accept", "application/json").
		SetContext(req.Context()).
		Get("https://jsonplaceholder.typicode.com/posts/" + id)

	if err != nil {
		logger.Error("Error during request to external API", "err", err)
		http.Error(writer, "External API error", http.StatusBadGateway)
		return
	}

	writer.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	writer.WriteHeader(resp.StatusCode())

	_, err = writer.Write(resp.Body())
	if err != nil {
		logger.Error("Error sending body to client", "err", err)
	}
}
