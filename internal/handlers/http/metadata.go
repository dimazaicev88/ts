package http

import (
	"context"

	"github.com/dimazaicev88/ts/internal/config"
	"github.com/dimazaicev88/ts/internal/dto"

	"github.com/gofiber/fiber/v3"
)

type Metadata struct {
	ctx          context.Context
	fb           *fiber.App
	serverConfig config.Config
}

func NewMetadata(ctx context.Context, fb *fiber.App, config config.Config) Metadata {
	return Metadata{
		fb:           fb,
		ctx:          ctx,
		serverConfig: config,
	}
}

// GetMetadata Добавить конфигурацию.
//
//	@Summary		Получить метаданные.
//	@Description	Получить метаданные.
//	@Tags			Server
//	@Produce		json
//	@Success		200	{object}	dto.ServerMetadata	"Метаданные сервера"
//	@Router			/api/v1/server/metadata [get]
func (m Metadata) GetMetadata(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(dto.ServerMetadata{
		NatsURL:      m.serverConfig.NatsURL,
		StreamName:   m.serverConfig.StreamName,
		ConsumerName: m.serverConfig.ConsumerName,
	})
}
