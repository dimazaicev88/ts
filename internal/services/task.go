package services

import (
	"context"
	"time"

	"github.com/dimazaicev88/ts/internal/converter"
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/storage"
	"github.com/dimazaicev88/ts/internal/storage/filters"
	"github.com/dimazaicev88/ts/internal/storage/models"
	"github.com/dimazaicev88/ts/internal/storage/repository"
	statusTask "github.com/dimazaicev88/ts/internal/task/status"
)

// ITask Работа с задачами
type ITask interface {
	Add(ctx context.Context, newTask dto.EventNewTask) error
	Update(ctx context.Context, uid string, fields map[string]string) error
	Delete(ctx context.Context, uids []string) error
	FindByUid(ctx context.Context, uid string) (dto.ResponseTask, error)
	Count(ctx context.Context, filter filters.Tasks) (int, error)
	Find(ctx context.Context, filter filters.Tasks) ([]dto.ResponseTask, error)
}

type Task struct {
	taskRepository repository.ITask
}

func NewTask(task repository.ITask) Task {
	return Task{taskRepository: task}
}

// Add Добавить новую задачу
func (t Task) Add(ctx context.Context, message dto.EventNewTask) error {
	taskModel := models.Task{
		Uid:        message.Uid,
		Status:     statusTask.Pending.ToString(),
		DateCreate: time.Now(),
		DateUpdate: time.Now(),
		Retention:  message.Retention,
		Subject:    message.Subject,
		WorkerUid:  message.WorkerUid,
		Timeout:    message.Timeout,
		Payload:    message.Payload,
	}

	return t.taskRepository.Add(ctx, taskModel)
}

// Update Обновить состояние задачи
func (t Task) Update(ctx context.Context, uid string, fields map[string]string) error {
	return t.taskRepository.Update(ctx, uid, fields)
}

// Delete Удалить задачу
func (t Task) Delete(ctx context.Context, uids []string) error {
	return t.taskRepository.DeleteByUid(ctx, uids)
}

// FindByUid Поиск задачи по uid
func (t Task) FindByUid(ctx context.Context, uid string) (dto.ResponseTask, error) {
	task, err := t.taskRepository.FindByUid(ctx, uid)
	if err != nil && !storage.IsErrNoRows(err) {
		return dto.ResponseTask{}, err
	}

	return converter.TaskModelToDto(task), nil
}

func (t Task) Count(ctx context.Context, filter filters.Tasks) (int, error) {
	return t.taskRepository.Count(ctx, filter)
}

func (t Task) Find(ctx context.Context, filter filters.Tasks) ([]dto.ResponseTask, error) {
	tasksModel, err := t.taskRepository.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return converter.ListTaskModelToListDto(tasksModel), err
}
