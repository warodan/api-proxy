package client

import (
	"api-proxy/internal/middleware"
	"github.com/go-resty/resty/v2"
)

func NewRestyClient() *resty.Client {
	client := resty.New()

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log := middleware.LoggerFromContext(resp.Request.Context())
		log.Info("resty request completed",
			"status", resp.StatusCode(),
			"url", resp.Request.URL,
			"duration", resp.Time(),
		)
		return nil
	})

	return client
}
