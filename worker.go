package ts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	config2 "github.com/dimazaicev88/ts/config"
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/services"
	statusTask "github.com/dimazaicev88/ts/internal/task/status"
	"github.com/dimazaicev88/ts/internal/wss"
	"github.com/dimazaicev88/ts/internal/wss/events"

	"github.com/dimazaicev88/ts/internal/base"

	"github.com/centrifugal/centrifuge-go"
	"github.com/nats-io/nats.go"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	config config2.WorkerConfig
	nc     *nats.Conn
	js     nats.JetStreamContext

	//Конфигурация сервера
	serverMetadata dto.ServerMetadata

	centrifugeClient *centrifuge.Client

	//Список обработчиков задач
	handlers          []base.HandlerConfig
	ctx               context.Context
	centrifugeService services.Centrifuge
	uid               string
}

func NewWorkerPool(
	ctx context.Context,
	config config2.WorkerConfig,
) (*Worker, error) {
	healService := services.NewHealthServer(config.ServerURL)
	err := healService.WaitServerAvailable(ctx, time.Hour*24, time.Second*10)
	if err != nil {
		return nil, err
	}

	serverAPI := services.NewServerAPI(config.ServerURL)
	serverMetadata, err := serverAPI.GetServerMetadata(ctx)
	if err != nil {
		return nil, err
	}

	nc, err := base.DefaultConnectToNats(serverMetadata.NatsURL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, err
	}

	centrifugeClient, err := wss.New()
	if err != nil {
		nc.Close()
		return nil, err
	}
	return &Worker{
		config:            config,
		nc:                nc,
		serverMetadata:    serverMetadata,
		js:                js,
		ctx:               ctx,
		centrifugeClient:  centrifugeClient,
		centrifugeService: services.NewCentrifuge(centrifugeClient),
		handlers:          make([]base.HandlerConfig, 0),
	}, nil
}

func (w *Worker) AddHandler(handlerConfig base.HandlerConfig) error {
	if handlerConfig.Handler == nil {
		return errors.New("handler is nil")
	}
	if handlerConfig.Subject == "" {
		return errors.New("subject is empty")
	}
	if handlerConfig.Concurrency <= 0 {
		return errors.New("concurrency must be greater than 0")
	}

	w.handlers = append(w.handlers, handlerConfig)
	return nil
}

func (w *Worker) Start() error {
	log.Print("starting workers")

	channel := make(chan string)
	subscribe, err := wss.CommonSubscribe(w.centrifugeClient)
	if err != nil {
		return err
	}

	//Ответ от сервера на запрос подключения
	subscribe.OnPublication(func(event centrifuge.PublicationEvent) {
		log.Debug().Msgf("centrifuge data: %s", event.Data)
		var eventData dto.EventData
		err := json.Unmarshal(event.Data, &eventData)
		if err != nil {
			log.Panic().Err(err).Msg("failed to unmarshal event data")
		}

		if eventData.EventName == events.AnswerWorkerConnect {
			var answer dto.EventAnswerWorkerConnect
			err = json.Unmarshal(eventData.Payload, &answer)
			if err != nil {
				log.Panic().Err(err).Msg("failed to unmarshal event data")
			}
			if answer.WorkerName == w.config.WorkerName {
				log.Debug().Msgf("getting worker uid: %s", answer.WorkerUid)
				channel <- answer.WorkerUid
			}
		}
	})

	hostname, _ := os.Hostname()
	w.centrifugeService.SendWorkerConnect(w.ctx, hostname, w.config.WorkerName)

	select {
	case w.uid = <-channel:
	case <-time.After(10 * time.Second):
		log.Panic().Msg("timeout waiting for worker")
	}

	g, _ := errgroup.WithContext(w.ctx)

	for _, handlerConfig := range w.handlers {
		h := handlerConfig
		g.Go(func() error {
			if err := w.work(w.ctx, h); err != nil {
				log.Print("worker error:", err)
				return err
			}
			return nil
		})
	}

	return g.Wait()
}

func (w *Worker) setDefaultValues(handlerConfig base.HandlerConfig) base.HandlerConfig {
	if handlerConfig.BatchSize == 0 {
		handlerConfig.BatchSize = 100
	}

	if handlerConfig.MaxWaitTime.Milliseconds() == 0 {
		handlerConfig.MaxWaitTime = time.Millisecond * 500
	}

	if handlerConfig.PollingTime.Milliseconds() == 0 {
		handlerConfig.PollingTime = time.Millisecond * 100
	}

	return handlerConfig
}

func (w *Worker) work(ctx context.Context, hc base.HandlerConfig) error {
	log.Info().
		Str("subject", hc.Subject).
		Msg("starting worker")

	handlerConfig := w.setDefaultValues(hc)

	sub, err := w.createSubscribe(handlerConfig)
	if err != nil {
		return err
	}

	pool, err := ants.NewPoolWithFunc(
		handlerConfig.Concurrency,
		w.makeHandler(ctx, handlerConfig),
		ants.WithPanicHandler(func(err interface{}) {
			log.Error().Interface("panic", err).Msg("panic in pool")
		}),
	)
	if err != nil {
		return err
	}
	defer pool.Release()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			log.Print("worker unsubscribing")
			sub.Unsubscribe()
			log.Print("worker unsubscribed")
			log.Print("worker pool releasing")
			pool.Release()
			log.Print("worker pool released")
			return nil

		case <-ticker.C:
			msgs, err := sub.Fetch(handlerConfig.BatchSize, nats.MaxWait(handlerConfig.MaxWaitTime))
			if err != nil {
				if !errors.Is(err, nats.ErrTimeout) {
					log.Error().Err(err).Msg("fetch error")
				}
				continue
			}

			for _, msg := range msgs {
				log.Print(len(msgs))
				err := pool.Invoke(base.TmpHandlerData{
					NatsMsg: msg,
					Data:    msg.Data,
					Ctx:     ctx,
				})
				if err != nil {
					log.Error().Err(err).Msg("pool invoke error")
				}
			}
		}
	}
}

func (w *Worker) createSubscribe(hc base.HandlerConfig) (*nats.Subscription, error) {
	sub, err := w.js.PullSubscribe(
		hc.Subject,
		w.serverMetadata.ConsumerName,
		nats.Bind(w.serverMetadata.StreamName, w.serverMetadata.ConsumerName),
	)

	if err != nil {
		return nil, fmt.Errorf("create PullSubscribe error: %w", err)
	}

	return sub, nil
}

func (w *Worker) makeHandler(ctx context.Context, hc base.HandlerConfig) func(data any) {
	return func(data any) {
		handlerData := data.(base.TmpHandlerData)
		var taskInfo base.TaskInfo
		err := json.Unmarshal(handlerData.Data, &taskInfo)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal task info")
		}

		// ACK сообщения
		if err := handlerData.NatsMsg.Ack(); err != nil {
			log.Error().Err(err).Msg("ack failed")
		}

		// Оборачиваем вызов Handler в recover
		func() {
			defer func() {
				if r := recover(); r != nil {
					errMsg := fmt.Sprintf("panic in handler: %v", r)
					log.Error().Str("task_uid", taskInfo.TaskUid).Msg(errMsg)

					// Обновляем статус задачи как ошибку
					updateErr := w.centrifugeService.SendUpdateTask(w.ctx, dto.EventUpdateTask{
						Uid: taskInfo.TaskUid,
						Fields: map[string]string{
							"error_msg": errMsg,
							"status":    statusTask.Archived.ToString(),
						},
					})
					if updateErr != nil {
						log.Error().Err(updateErr).Msg("handle update error after panic")
					}
				}
			}()

			err = hc.Handler(ctx, base.HandlerData{
				NatsMsg:  handlerData.NatsMsg,
				Ctx:      handlerData.Ctx,
				TaskInfo: taskInfo,
			})
		}()

		// Обрабатываем результат
		if err != nil {
			updateErr := w.centrifugeService.SendUpdateTask(w.ctx, dto.EventUpdateTask{
				Uid: taskInfo.TaskUid,
				Fields: map[string]string{
					"errorMsg": err.Error(),
					"status":   statusTask.Archived.ToString(),
				},
			})
			if updateErr != nil {
				log.Error().Err(updateErr).Msg("handle update error")
			}
		} else {
			updateErr := w.centrifugeService.SendUpdateTask(w.ctx, dto.EventUpdateTask{
				Uid: taskInfo.TaskUid,
				Fields: map[string]string{
					"status": statusTask.Completed.ToString(),
				},
			})
			if updateErr != nil {
				log.Error().Err(updateErr).Msg("handle update error")
			}
		}
	}
}

func (w *Worker) Stop() {
	log.Info().Msg("stopping worker")

	if w.nc != nil {
		_ = w.nc.Drain() // graceful
		w.nc.Close()
	}
	log.Info().Msg("workers stopped")
}
