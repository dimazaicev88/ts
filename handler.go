package ts

import (
	"context"

	"github.com/nats-io/nats.go"
)

type HandlerData struct {
	NatsMsg  *nats.Msg
	Ctx      context.Context
	TaskInfo TaskInfo
}

type HandlerSubject func(ctx context.Context, handlerData HandlerData) error

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
