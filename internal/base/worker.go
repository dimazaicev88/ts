package base

import (
	"context"

	"github.com/nats-io/nats.go"
)

type WorkerStatus int8

const Disconnected WorkerStatus = 1
const Connected WorkerStatus = 1

func (w WorkerStatus) ToInt8() int8 {
	return int8(w)
}

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

type HandlerData struct {
	NatsMsg  *nats.Msg
	Ctx      context.Context
	TaskInfo TaskInfo
}

type HandlerSubject func(ctx context.Context, handlerData HandlerData) error
