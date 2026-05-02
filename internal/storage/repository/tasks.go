package repository

import (
	"context"
	"errors"

	"github.com/dimazaicev88/ts/internal/storage"
	"github.com/dimazaicev88/ts/internal/storage/filters"
	"github.com/dimazaicev88/ts/internal/storage/models"

	"github.com/uptrace/bun"
)

type ITask interface {
	Add(ctx context.Context, task models.Task) error
	Update(ctx context.Context, uid string, fields map[string]string) error
	DeleteByUid(ctx context.Context, uids []string) error
	FindByUid(ctx context.Context, uid string) (models.Task, error)
	Count(ctx context.Context, filter filters.Tasks) (int, error)
	Find(ctx context.Context, filter filters.Tasks) ([]models.Task, error)
}

type Task struct {
	db *bun.DB
}

func NewTask(db *bun.DB) Task {
	return Task{db: db}
}

func (t Task) Add(ctx context.Context, task models.Task) error {
	_, err := t.db.NewInsert().Model(&task).Exec(ctx)
	return err
}

func (t Task) Update(ctx context.Context, uid string, fields map[string]string) error {
	if len(fields) == 0 {
		return errors.New("fields is empty")
	}

	update := t.db.NewUpdate().Model(new(models.Task))
	for name, value := range fields {
		update.Set(name+"=?", value)
	}
	_, err := update.Where("uid=?", uid).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t Task) DeleteByUid(ctx context.Context, uids []string) error {
	_, err := t.db.NewDelete().Model(new(models.Task)).Where("uid IN (?)", bun.List(uids)).Exec(ctx)
	return err
}

func (t Task) FindByUid(ctx context.Context, uid string) (models.Task, error) {
	var taskModel models.Task
	err := t.db.NewSelect().Model(&taskModel).Where("uid = ?", uid).Scan(ctx)
	if err != nil {
		return models.Task{}, err
	}
	return taskModel, nil
}

func (t Task) Count(ctx context.Context, filter filters.Tasks) (int, error) {
	if len(filter.WorkerUid) == 0 {
		return 0, errors.New("filter worker uid is empty")
	}

	query := t.db.NewSelect().Model(new(models.Task))
	query.Where("worker_uid=?", filter.WorkerUid)
	if len(filter.Status) > 0 {
		query.Where("status=?", filter.Status)
	}
	return query.Count(ctx)
}

func (t Task) Find(ctx context.Context, filter filters.Tasks) ([]models.Task, error) {
	if len(filter.WorkerUid) == 0 {
		return nil, errors.New("filter worker uid is empty")
	}

	var tasks []models.Task
	query := t.db.NewSelect().
		Model(&tasks)
	query.Limit(filter.Limit)
	if filter.Skip > 0 {
		query.Offset(filter.Skip)
	}
	query.Where("worker_uid=?", filter.WorkerUid)

	if len(filter.Status) > 0 {
		query.Where("status=?", filter.Status)
	}

	if err := query.Scan(ctx, &tasks); err != nil && storage.IsNoErrNoRows(err) {
		return nil, err
	}

	return tasks, nil
}
