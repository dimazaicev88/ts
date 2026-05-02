package services

import (
	"context"
	"net/http"
	"time"

	"github.com/dimazaicev88/ts/internal/tools"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type HealthServer struct {
	serverURL string
}

func NewHealthServer(serverURL string) HealthServer {
	return HealthServer{serverURL: serverURL}
}

func (h HealthServer) WaitServerAvailable(ctx context.Context, timeout time.Duration, interval time.Duration) error {
	restyClient := resty.New()
	restyClient.SetTimeout(5 * time.Second).SetBaseURL(h.serverURL)

	return tools.Until(ctx, timeout, interval, "time out", func() (completed bool, err error) {
		log.Debug().Msg("Waiting server available")
		head, _ := restyClient.R().SetContext(ctx).Head("/api/v1/server/metadata")
		return head.StatusCode() == http.StatusOK, nil
	})
}
