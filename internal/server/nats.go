package server

import (
	"context"
	"errors"
	"time"

	"github.com/dimazaicev88/ts/internal/config"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

// CreateStreamWithConsumer Создаём Stream
func CreateStreamWithConsumer(ctx context.Context, config config.Config) error {
	nc, err := nats.Connect(
		config.NatsURL,
		nats.Timeout(5*time.Second),
		nats.ReconnectWait(5*time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			log.Error().Err(err).Msg("NATS disconnected")
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Debug().Str("url", c.ConnectedUrl()).Msg("NATS reconnected")
		}),
	)
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}

	stream, err := js.Stream(ctx, config.StreamName)
	if err != nil {
		var jsErr jetstream.JetStreamError
		ok := errors.As(err, &jsErr)
		if ok && jsErr.APIError().ErrorCode != jetstream.JSErrCodeStreamNotFound {
			return err
		}
	}

	if stream == nil {
		stream, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      config.StreamName,
			Subjects:  []string{config.SubjectName},
			Storage:   jetstream.FileStorage,
			Retention: jetstream.WorkQueuePolicy, // Для очереди задач
			MaxAge:    744 * time.Hour,           //744h = 31d * 24h
		})
		if err != nil {
			return err
		}

		log.Debug().Msgf("Stream '%s' created", config.StreamName)

	} else {
		log.Debug().Msgf("Stream '%s' already exists", config.StreamName)
	}
	_, err = stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       config.ConsumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxAckPending: -1,
		MaxDeliver:    3, // Максимум 3 попытки доставки
	})

	if err != nil {
		return err
	}

	log.Debug().Msg("Consumer created")

	return nil
}
