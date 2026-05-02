package services

import (
	"context"
	"encoding/json"

	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/tools"
	"github.com/dimazaicev88/ts/internal/wss"
	"github.com/dimazaicev88/ts/internal/wss/events"

	"github.com/centrifugal/centrifuge-go"
	"github.com/rs/zerolog/log"
)

type Centrifuge struct {
	centrifugeClient *centrifuge.Client
}

func NewCentrifuge(centrifugeClient *centrifuge.Client) Centrifuge {
	return Centrifuge{
		centrifugeClient: centrifugeClient,
	}
}

func (c Centrifuge) SendAnswerWorkerConnect(
	ctx context.Context,
	workerUid string,
	workerName string,
) {
	data := dto.EventData{
		EventName: events.AnswerWorkerConnect,
		Payload: tools.ToJson(
			dto.EventAnswerWorkerConnect{
				WorkerUid:  workerUid,
				WorkerName: workerName,
			},
		),
	}
	eventData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling data")
	}
	_, err = c.centrifugeClient.Publish(ctx, wss.CommonChannel, eventData)
	if err != nil {
		log.Error().Err(err).Msg("error publishing event")
	}
}

func (c Centrifuge) SendWorkerConnect(
	ctx context.Context,
	hostname string,
	workerName string,
) {
	data := dto.EventData{
		EventName: events.WorkerConnect,
		Payload: tools.ToJson(
			dto.EventWorkerConnect{
				WorkerName: workerName,
				Hostname:   hostname,
			},
		),
	}

	eventData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling data")
	}
	_, err = c.centrifugeClient.Publish(ctx, wss.CommonChannel, eventData)
	if err != nil {
		log.Error().Err(err).Msg("error publishing event")
	}
}

func (c Centrifuge) SendNewTask(
	ctx context.Context,
	newTask dto.EventNewTask,
) error {
	data := dto.EventData{
		EventName: events.NewTaskConnect,
		Payload:   tools.ToJson(newTask),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.centrifugeClient.Publish(ctx, wss.CommonChannel, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (c Centrifuge) SendUpdateTask(
	ctx context.Context,
	updateTask dto.EventUpdateTask,
) error {
	data := dto.EventData{
		EventName: events.UpdateTask,
		Payload:   tools.ToJson(updateTask),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.centrifugeClient.Publish(ctx, wss.CommonChannel, jsonData)
	if err != nil {
		return err
	}

	return nil
}
