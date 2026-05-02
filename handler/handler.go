package handler

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Data struct {
	NatsMsg  *nats.Msg
	Ctx      context.Context
	TaskInfo TaskInfo
}

type Subject func(ctx context.Context, handlerData Data) error

type TmpHandlerData struct {
	NatsMsg *nats.Msg
	Data    []byte
	Ctx     context.Context
}

type TaskInfo struct {
	TaskUid   string `json:"taskUid"`
	Payload   []byte `json:"payload"`
	Timeout   uint64 `json:"timeout"`
	WorkerUid string `json:"workerUid"`
}
