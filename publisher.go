package ts

import (
	"context"
	"encoding/json"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/dimazaicev88/ts/internal/base"
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/services"
	"github.com/dimazaicev88/ts/internal/wss"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

type Publisher struct {
	nc                *nats.Conn
	js                jetstream.JetStream
	streamName        string
	centrifugeClient  *centrifuge.Client
	centrifugeService services.Centrifuge
	workerUid         string
}

func NewPublisher(
	ctx context.Context,
	workerName string,
	serverURL string,
	timeoutWaitServer time.Duration,
) (*Publisher, error) {
	healService := services.NewHealthServer(serverURL)
	err := healService.WaitServerAvailable(ctx, timeoutWaitServer, 10*time.Second)
	if err != nil {
		return nil, err
	}

	serverAPI := services.NewServerAPI(serverURL)
	serverMetadata, err := serverAPI.GetServerMetadata(ctx)
	if err != nil {
		return nil, err
	}

	worker, err := serverAPI.FindWorkerByName(ctx, workerName)
	if err != nil {
		return nil, err
	}

	nc, err := base.DefaultConnectToNats(serverMetadata.NatsURL)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, err
	}

	// Проверяем доступность JetStream
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := js.AccountInfo(ctx); err != nil {
		nc.Close()
		return nil, err
	}

	centrifugeClient, err := wss.New()
	if err != nil {
		nc.Close()
		return nil, err
	}
	return &Publisher{
		nc:                nc,
		js:                js,
		streamName:        serverMetadata.StreamName,
		centrifugeClient:  centrifugeClient,
		centrifugeService: services.NewCentrifuge(centrifugeClient),
		workerUid:         worker.Uid,
	}, nil
}

func (p *Publisher) Publish(
	ctx context.Context,
	subjectName string,
	payload string,
	timeout uint64,
	retention uint64,
) (string, error) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	newTask := dto.EventNewTask{
		Uid:       uuid.New().String(),
		Payload:   payload,
		Subject:   subjectName,
		Timeout:   timeout,
		Retention: retention,
		WorkerUid: p.workerUid,
	}

	err = p.centrifugeService.SendNewTask(ctx, newTask)
	if err != nil {
		return "", err
	}

	// Публикуем с ожиданием подтверждения
	ack, err := p.js.Publish(ctx, subjectName, jsonData)
	if err != nil {
		// Логируем детали ошибки
		log.Error().
			Err(err).
			Str("subject", subjectName).
			Str("stream", p.streamName).
			Msg("Failed to publish message")
		return "", err
	}

	log.Debug().
		Uint64("sequence", ack.Sequence).
		Str("subject", subjectName).
		Msg("Message published successfully")

	return newTask.Uid, nil
}

func (p *Publisher) PublishAsync(subjectName string, data []byte) (jetstream.PubAckFuture, error) {
	return p.js.PublishAsync(subjectName, data)
}

func (p *Publisher) Close() {
	p.centrifugeClient.Close()
	p.nc.Close()
}
