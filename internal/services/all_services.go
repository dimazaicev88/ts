package services

import (
	"github.com/dimazaicev88/ts/internal/storage/repository"

	"github.com/centrifugal/centrifuge-go"
	"github.com/uptrace/bun"
)

type AllServices struct {
	TaskService   ITask
	WorkerService IWorker
	Centrifuge    Centrifuge
}

func NewAllServices(db *bun.DB, centrifugeClient *centrifuge.Client) AllServices {
	return AllServices{
		TaskService:   NewTask(repository.NewTask(db)),
		WorkerService: NewWorker(repository.NewWorker(db)),
		Centrifuge:    NewCentrifuge(centrifugeClient),
	}
}
