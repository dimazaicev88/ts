package services

import (
	"context"
	"testing"
	"time"

	"github.com/dimazaicev88/ts/internal/base"
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/fixtures"
	"github.com/dimazaicev88/ts/internal/storage"
	"github.com/dimazaicev88/ts/internal/storage/models"
	"github.com/dimazaicev88/ts/internal/storage/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWorker(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	ctx := context.Background()
	dbName := "tusk_runner"
	testContainer := fixtures.NewMysqlTestContainer(ctx, req, dbName)
	db, err := fixtures.SetupTestContainerMysql(ctx, testContainer, "/home/dima/projects/task-runner/migrations/", dbName)
	req.NoError(err)
	workerRepository := repository.NewWorker(db)
	workerService := NewWorker(workerRepository)

	t.Run("Add", func(t *testing.T) {
		storage.TruncateAllTables(ctx, db)
		dtoWorker := dto.Worker{
			Uid:    uuid.NewString(),
			Name:   "test-worker",
			Status: 1,
		}
		err = workerService.Add(ctx, dtoWorker)
		req.NoError(err)

		var worker models.Worker
		dbQuery := db.NewSelect().Model(&worker).Where("uid = ?", dtoWorker.Uid)

		err := dbQuery.Scan(ctx)
		req.NoError(err)

		req.NoError(err)
		req.Equal(dtoWorker.Uid, worker.Uid)
		req.Equal(dtoWorker.Name, worker.Name)
		req.Equal(int8(1), worker.Status)
	})

	t.Run("FindByUid", func(t *testing.T) {
		storage.TruncateAllTables(ctx, db)
		model := models.Worker{
			Uid:           uuid.NewString(),
			Name:          "test-worker",
			LastConnected: time.Now(),
			Status:        int8(base.Connected),
		}
		_, err = db.NewInsert().Model(&model).Exec(ctx)
		req.NoError(err)

		worker, err := workerService.FindByUid(ctx, model.Uid)
		req.NoError(err)

		req.NoError(err)
		req.Equal(model.Uid, worker.Uid)
		req.Equal(model.Name, worker.Name)
		req.Equal(int8(1), worker.Status)
	})

	t.Run("FindByName", func(t *testing.T) {
		storage.TruncateAllTables(ctx, db)
		model := models.Worker{
			Uid:           uuid.NewString(),
			Name:          "test-worker",
			LastConnected: time.Now(),
			Status:        int8(base.Connected),
		}
		_, err = db.NewInsert().Model(&model).Exec(ctx)
		req.NoError(err)

		worker, err := workerService.FindByName(ctx, model.Name)
		req.NoError(err)

		req.NoError(err)
		req.Equal(model.Uid, worker.Uid)
		req.Equal(model.Name, worker.Name)
		req.Equal(int8(1), worker.Status)
	})

	t.Run("DeleteByUid", func(t *testing.T) {
		storage.TruncateAllTables(ctx, db)
		model := models.Worker{
			Uid:           uuid.NewString(),
			Name:          "test-worker",
			LastConnected: time.Now(),
			Status:        int8(base.Connected),
		}
		_, err = db.NewInsert().Model(&model).Exec(ctx)
		req.NoError(err)

		err := workerService.Delete(ctx, model.Uid)
		req.NoError(err)

		var worker models.Worker
		cnt, err := db.NewSelect().Model(&worker).Where("uid = ?", model.Uid).Count(ctx)
		req.NoError(err)
		req.Equal(cnt, 0)
	})
}
