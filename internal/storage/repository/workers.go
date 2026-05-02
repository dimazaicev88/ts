package repository

import (
	"context"
	"time"

	"github.com/dimazaicev88/ts/internal/storage/filters"
	"github.com/dimazaicev88/ts/internal/storage/models"

	"github.com/uptrace/bun"
)

type IWorker interface {
	Add(ctx context.Context, task models.Worker) error
	Update(ctx context.Context) error
	DeleteByUid(ctx context.Context, uid string) error
	FindByUid(ctx context.Context, uid string) (models.Worker, error)
	FindByName(ctx context.Context, uid string) (models.Worker, error)
	Heartbeat(ctx context.Context, uid string) error
	Find(ctx context.Context, filter filters.Worker) ([]models.Worker, error)
	Count(ctx context.Context, filter filters.Worker) (int, error)
}

type Worker struct {
	db *bun.DB
}

func NewWorker(db *bun.DB) Worker {
	return Worker{db: db}
}

func (w Worker) Add(ctx context.Context, task models.Worker) error {
	_, err := w.db.NewInsert().Model(&task).Exec(ctx)
	return err
}

func (w Worker) Update(ctx context.Context) error {
	return nil
}

func (w Worker) DeleteByUid(ctx context.Context, uid string) error {
	_, err := w.db.NewDelete().Model(new(models.Worker)).Where("uid = ?", uid).Exec(ctx)
	return err
}

func (w Worker) FindByUid(ctx context.Context, uid string) (models.Worker, error) {
	var taskModel models.Worker
	dbQuery := w.db.NewSelect().Model(&taskModel).Where("uid = ?", uid)

	err := dbQuery.Scan(ctx)
	if err != nil {
		return models.Worker{}, err
	}
	return taskModel, nil
}

func (w Worker) FindByName(ctx context.Context, uid string) (models.Worker, error) {
	var worker models.Worker
	dbQuery := w.db.NewSelect().Model(&worker).Where("name = ?", uid)

	err := dbQuery.Scan(ctx)
	if err != nil {
		return models.Worker{}, err
	}
	return worker, nil
}

func (w Worker) Heartbeat(ctx context.Context, uid string) error {
	_, err := w.db.NewUpdate().Model(new(models.Worker)).
		Set("last_connected=?", time.Now()).
		Where("uid=?", uid).Exec(ctx)
	return err
}

func (w Worker) Find(ctx context.Context, filter filters.Worker) ([]models.Worker, error) {
	var workers []models.Worker
	dbQuery := w.db.NewSelect().Model(&workers)

	if filter.Name.Valid {
		dbQuery.Where("name LIKE CONCAT('%', ?, '%')", filter.Name.V)
	}

	if filter.Status.Valid {
		dbQuery.Where("status=?", filter.Status.V)
	}

	dbQuery.Limit(filter.Limit)

	if filter.Skip > 0 {
		dbQuery.Offset(filter.Skip)
	}

	err := dbQuery.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return workers, nil
}

func (w Worker) Count(ctx context.Context, filter filters.Worker) (int, error) {
	dbQuery := w.db.NewSelect().Model(&models.Worker{})

	if filter.Name.Valid {
		dbQuery.Where("name LIKE CONCAT('%', ?, '%')", filter.Name.V)
	}

	if filter.Status.Valid {
		dbQuery.Where("status=?", filter.Status.V)
	}

	return dbQuery.Count(ctx)
}
