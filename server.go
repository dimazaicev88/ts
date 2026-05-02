package ts

import (
	"context"
	"fmt"

	"github.com/dimazaicev88/ts/config"
	_ "github.com/dimazaicev88/ts/docs"
	wssHandlers "github.com/dimazaicev88/ts/internal/handlers/ws"
	"github.com/dimazaicev88/ts/internal/services"
	"github.com/dimazaicev88/ts/internal/storage"
	"github.com/dimazaicev88/ts/internal/wss"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/rs/zerolog/log"
)

func RunServer(ctx context.Context, config config.ServerConfig) {
	if config.HttpServerPort <= 0 {
		log.Fatal().Msg("http port must be greater than zero")
	}

	db, err := storage.NewMysql(config.DbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to db failed")
	}

	//Fiber init
	fb := fiber.New()
	fb.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-KEY"},
	}))

	centrifugeClient, err := wss.New()
	if err != nil {
		log.Fatal().Err(err).Msg("connect to centrifugo client failed")
	}
	allServices := services.NewAllServices(db, centrifugeClient)
	NewRoutes(ctx, fb, config, allServices).AddRoutes()
	wssHandlers.RegWssHandlers(ctx, allServices, centrifugeClient)
	errStartServer := make(chan error)

	//Запуск fiber
	go func() {
		log.Debug().Msg("starting fiber server")
		errStartServer <- fb.Listen(fmt.Sprintf(":%d", config.HttpServerPort))
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("shutdown server")
		centrifugeClient.Close()
		db.Close()
	case err := <-errStartServer:
		db.Close()
		centrifugeClient.Close()
		log.Fatal().Err(err).Msg("server error")
	}
}
