package services

import (
	"context"
	"time"

	"github.com/dimazaicev88/ts/internal/dto"

	"github.com/go-resty/resty/v2"
	"github.com/segmentio/encoding/json"
)

type ServerAPI struct {
	serverURL string
}

func NewServerAPI(serverURL string) *ServerAPI {
	return &ServerAPI{serverURL: serverURL}
}

func (s ServerAPI) ServerMetadata(timeout time.Duration) (dto.ServerMetadata, error) {
	restyClient := resty.New()
	restyClient.SetTimeout(timeout).SetBaseURL(s.serverURL)

	data, err := restyClient.R().Get("/api/v1/server/metadata")
	if err != nil {
		return dto.ServerMetadata{}, err
	}

	var serverMetadata dto.ServerMetadata
	err = json.Unmarshal(data.Body(), &serverMetadata)
	if err != nil {
		return dto.ServerMetadata{}, err
	}

	return serverMetadata, nil
}

func (s ServerAPI) GetServerMetadata(ctx context.Context) (dto.ServerMetadata, error) {
	restyClient := resty.New()
	restyClient.SetTimeout(time.Second * 30).SetBaseURL(s.serverURL)

	data, err := restyClient.R().SetContext(ctx).Get("/api/v1/server/metadata")
	if err != nil {
		return dto.ServerMetadata{}, err
	}

	var serverMetadata dto.ServerMetadata
	err = json.Unmarshal(data.Body(), &serverMetadata)
	if err != nil {
		return dto.ServerMetadata{}, err
	}

	return serverMetadata, nil
}

func (s ServerAPI) FindWorkerByName(ctx context.Context, name string) (dto.Worker, error) {
	restyClient := resty.New()
	restyClient.SetTimeout(time.Second * 30).SetBaseURL(s.serverURL)

	data, err := restyClient.R().SetContext(ctx).Get("/api/v1/wroker/" + name)
	if err != nil {
		return dto.Worker{}, err
	}

	var worker dto.Worker
	err = json.Unmarshal(data.Body(), &worker)
	if err != nil {
		return dto.Worker{}, err
	}

	return worker, nil
}
