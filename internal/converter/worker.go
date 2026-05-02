package converter

import (
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/storage/models"
)

func WorkerModelToDto(models models.Worker) dto.Worker {
	return dto.Worker{
		Uid:           models.Uid,
		Name:          models.Name,
		LastConnected: models.LastConnected,
		Status:        models.Status,
	}
}

func ListWorkerModelToDto(models []models.Worker) []dto.Worker {
	var workers []dto.Worker
	for _, model := range models {
		workers = append(workers, WorkerModelToDto(model))
	}

	return workers
}
