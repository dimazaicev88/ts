package services

import (
	"context"

	"github.com/dimazaicev88/ts/internal/base"
	"github.com/dimazaicev88/ts/internal/converter"
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/storage"
	"github.com/dimazaicev88/ts/internal/storage/filters"
	"github.com/dimazaicev88/ts/internal/storage/models"
	"github.com/dimazaicev88/ts/internal/storage/repository"
)

// IWorker Работа с воркерами
type IWorker interface {
	Add(ctx context.Context, message dto.Worker) error
	Update(ctx context.Context) error
	Delete(ctx context.Context, uid string) error
	FindByUid(ctx context.Context, uid string) (dto.Worker, error)
	FindByName(ctx context.Context, name string) (dto.Worker, error)
	Find(ctx context.Context, filter filters.Worker) ([]dto.Worker, error)
	Heartbeat(ctx context.Context, uid string) error
	Count(ctx context.Context, filter filters.Worker) (int, error)
}
type Worker struct {
	workerRepository repository.IWorker
}

func NewWorker(workerRepository repository.IWorker) *Worker {
	return &Worker{
		workerRepository: workerRepository,
	}
}

// Add Добавить воркера
func (w *Worker) Add(ctx context.Context, worker dto.Worker) error {
	return w.workerRepository.Add(ctx, models.Worker{
		Uid:           worker.Uid,
		Name:          worker.Name,
		LastConnected: worker.LastConnected,
		Status:        int8(base.Connected),
	})
}

// Update Обновить воркера
func (w *Worker) Update(ctx context.Context) error {
	return w.workerRepository.Update(ctx)
}

func (w *Worker) Delete(ctx context.Context, uid string) error {
	return w.workerRepository.DeleteByUid(ctx, uid)
}

func (w *Worker) FindByUid(ctx context.Context, uid string) (dto.Worker, error) {
	workerModel, err := w.workerRepository.FindByUid(ctx, uid)
	if err != nil && !storage.IsErrNoRows(err) {
		return dto.Worker{}, err
	}

	return converter.WorkerModelToDto(workerModel), nil
}

func (w *Worker) FindByName(ctx context.Context, name string) (dto.Worker, error) {
	workerModel, err := w.workerRepository.FindByName(ctx, name)
	if err != nil && !storage.IsErrNoRows(err) {
		return dto.Worker{}, err
	}

	return converter.WorkerModelToDto(workerModel), nil
}

func (w *Worker) Heartbeat(ctx context.Context, uid string) error {
	return w.workerRepository.Heartbeat(ctx, uid)
}

func (w *Worker) Find(ctx context.Context, filter filters.Worker) ([]dto.Worker, error) {
	workers, err := w.workerRepository.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return converter.ListWorkerModelToDto(workers), nil
}

func (w *Worker) Count(ctx context.Context, filter filters.Worker) (int, error) {
	return w.workerRepository.Count(ctx, filter)
}
