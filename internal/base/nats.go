package base

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func DefaultConnectToNats(natsURL string) (*nats.Conn, error) {
	return nats.Connect(
		natsURL,
		nats.Timeout(5*time.Second),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			log.Error().Err(err).Msg("NATS disconnected")
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Debug().Str("url", c.ConnectedUrl()).Msg("NATS reconnected")
		}),
	)
}
