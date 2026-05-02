package ws

import (
	"context"
	"time"

	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/services"
	"github.com/dimazaicev88/ts/internal/wss"
	"github.com/dimazaicev88/ts/internal/wss/events"

	"github.com/centrifugal/centrifuge-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/encoding/json"
)

func handlerWorkerConnect(
	ctx context.Context,
	data []byte,
	allServices services.AllServices,
) {
	var workerConnect dto.EventWorkerConnect
	err := json.Unmarshal(data, &workerConnect)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling newTask")
		return
	}

	worker, err := allServices.WorkerService.FindByName(ctx, workerConnect.WorkerName)
	if err != nil {
		log.Error().Err(err).Msg("error finding worker")
	}

	// Если воркер уже существует - возвращаем его UID
	if !worker.IsEmpty() {
		allServices.Centrifuge.SendAnswerWorkerConnect(ctx, worker.Uid, worker.Name)
	} else {
		workerUid := uuid.NewString()
		err := allServices.WorkerService.Add(ctx, dto.Worker{
			Uid:           workerUid,
			Name:          workerConnect.WorkerName,
			LastConnected: time.Now(),
			Status:        1,
		})
		if err != nil {
			log.Error().Err(err).Msg("error adding worker")
		} else {
			allServices.Centrifuge.SendAnswerWorkerConnect(ctx, workerUid, workerConnect.WorkerName)
		}
	}
}

func handlerUpdateTask(
	ctx context.Context,
	data []byte,
	allServices services.AllServices,
) {
	var newTask dto.EventNewTask
	err := json.Unmarshal(data, &newTask)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling newTask")
		return
	}

	err = allServices.TaskService.Add(ctx, newTask)
	if err != nil {
		log.Error().Err(err).Msg("error finding worker")
		return
	}
}

func RegWssHandlers(
	ctx context.Context,
	allServices services.AllServices,
	centrifugeClient *centrifuge.Client,
) {
	subscribe, err := wss.CommonSubscribe(centrifugeClient)
	if err != nil {
		log.Panic().Err(err).Msg("error subscribing to common channel")
	}
	subscribe.OnPublication(func(e centrifuge.PublicationEvent) {
		var eventConnect dto.EventData

		err := json.Unmarshal(e.Data, &eventConnect)
		if err != nil {
			log.Error().Err(err).Msg("error unmarshalling event command")
			return
		}

		if eventConnect.EventName == events.WorkerConnect {
			handlerWorkerConnect(ctx, eventConnect.Payload, allServices)
		}

		if eventConnect.EventName == events.NewTaskConnect {
			handlerUpdateTask(ctx, eventConnect.Payload, allServices)
		}
	})
}
