package wss

import (
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/rs/zerolog/log"
)

// const CentrifugeURL = "ws://centrifugo.intsite.org/connection/websocket"
const CentrifugeURL = "ws://localhost:8100/connection/websocket"

func New() (*centrifuge.Client, error) {
	client := centrifuge.NewJsonClient(CentrifugeURL, centrifuge.Config{
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	})

	client.OnError(func(e centrifuge.ErrorEvent) {
		log.Error().Err(e.Error).Msg("centrifuge error")
	})

	client.OnConnected(func(e centrifuge.ConnectedEvent) {
		log.Info().Msg("connected to centrifuge")
	})

	client.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Warn().
			Str("reason", e.Reason).
			Msg("centrifuge disconnected")
	})

	err := client.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CommonSubscribe(client *centrifuge.Client) (*centrifuge.Subscription, error) {
	var err error
	sub, _ := client.GetSubscription(CommonChannel)
	if sub == nil {
		sub, err = client.NewSubscription(CommonChannel)
		if err != nil {
			return nil, err
		}
	}

	sub.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Error().Err(e.Error).Msg("subscription error")
	})

	if err := sub.Subscribe(); err != nil {
		return nil, err
	}

	return sub, nil
}
