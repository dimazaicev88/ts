package services

//
//import (
//	"context"
//	"github.com/dimazaicev88/ts/internal/base"
//	"github.com/dimazaicev88/ts/internal/fixtures"
//	"github.com/dimazaicev88/ts/internal/storage"
//	"github.com/dimazaicev88/ts/internal/storage/models"
//	"github.com/dimazaicev88/ts/internal/storage/repository"
//	statusTask "task-runner/internal/task/status"
//	"testing"
//	"time"
//
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/require"
//)
//
//func TestTask(t *testing.T) {
//	t.Parallel()
//	req := require.New(t)
//	ctx := context.Background()
//	dbName := "tusk_runner"
//	testContainer := fixtures.NewMysqlTestContainer(ctx, req, dbName)
//	db, err := fixtures.SetupTestContainerMysql(ctx, testContainer, "/home/dima/projects/task-runner/migrations/", dbName)
//	//db, err := fixtures.SetupByUrlConnect(
//	//	ctx,
//	//	"root:rootpassword@tcp(localhost:3306)/task_runner?multiStatements=true",
//	//	"/home/dima/projects/task-runner/migrations/",
//	//	dbName,
//	//)
//	req.NoError(err)
//	taskRepository := repository.NewTask(db)
//	taskService := NewTask(taskRepository)
//
//	t.Run("Add", func(t *testing.T) {
//		storage.TruncateAllTables(ctx, db)
//		modelWorker := models.Worker{
//			Uid:           uuid.NewString(),
//			Name:          "test-worker",
//			LastConnected: time.Now(),
//			Status:        int8(base.Connected),
//		}
//		_, err = db.NewInsert().Model(&modelWorker).Exec(ctx)
//		req.NoError(err)
//
//		taskMessage := grpcimpl.RequestAddTaskTaskMessage{
//			Payload:   []byte("test payload"),
//			Uid:       uuid.New().String(),
//			Subject:   "kor",
//			ErrorMsg:  "error message",
//			Timeout:   100,
//			Retention: 10,
//			WorkerUid: modelWorker.Uid,
//		}
//		err = taskService.Add(ctx, &taskMessage)
//		req.NoError(err)
//
//		var task models.Task
//		dbQuery := db.NewSelect().Model(&task).Where("uid = ?", taskMessage.Uid)
//
//		err := dbQuery.Scan(ctx)
//		req.NoError(err)
//
//		req.Equal(string(taskMessage.Payload), task.Payload)
//		req.Equal(taskMessage.Uid, task.Uid)
//		req.Equal(taskMessage.Subject, task.Subject)
//		req.Equal(taskMessage.ErrorMsg, task.ErrorMsg)
//		req.Equal(taskMessage.Timeout, task.Timeout)
//		req.Equal(taskMessage.Retention, task.Retention)
//		req.Equal(statusTask.Active.ToString(), task.Status)
//		req.Equal(taskMessage.WorkerUid, task.WorkerUid)
//	})
//
//	t.Run("FindByUid", func(t *testing.T) {
//		storage.TruncateAllTables(ctx, db)
//		modelWorker := models.Worker{
//			Uid:           uuid.NewString(),
//			Name:          "test-worker",
//			LastConnected: time.Now(),
//			Status:        int8(base.Connected),
//		}
//		_, err = db.NewInsert().Model(&modelWorker).Exec(ctx)
//		req.NoError(err)
//
//		taskMessage :=
//			models.Task{
//				Payload:     "test payload",
//				Uid:         uuid.New().String(),
//				Subject:     "kor",
//				ErrorMsg:    "error message",
//				Timeout:     100,
//				Retention:   10,
//				CompletedAt: 500,
//				Result:      "empty result",
//				Status:      "running",
//				WorkerUid:   modelWorker.Uid,
//			}
//		_, err = db.NewInsert().Model(&taskMessage).Exec(ctx)
//		req.NoError(err)
//
//		task, err := taskService.FindByUid(ctx, taskMessage.Uid)
//		req.NoError(err)
//
//		req.Equal(string(taskMessage.Payload), task.Payload)
//		req.Equal(taskMessage.Uid, task.Uid)
//		req.Equal(taskMessage.Subject, task.Subject)
//		req.Equal(taskMessage.ErrorMsg, task.ErrorMsg)
//		req.Equal(taskMessage.Timeout, task.Timeout)
//		req.Equal(taskMessage.Retention, task.Retention)
//		req.Equal(taskMessage.CompletedAt, task.CompletedAt)
//		req.Equal(taskMessage.Status, task.Status)
//		req.Equal(taskMessage.WorkerUid, task.WorkerUid)
//	})
//
//	t.Run("DeleteByUid", func(t *testing.T) {
//		storage.TruncateAllTables(ctx, db)
//		modelWorker := models.Worker{
//			Uid:           uuid.NewString(),
//			Name:          "test-worker",
//			LastConnected: time.Now(),
//			Status:        int8(base.Connected),
//		}
//		_, err = db.NewInsert().Model(&modelWorker).Exec(ctx)
//		req.NoError(err)
//
//		taskMessage :=
//			models.Task{
//				Payload:     "test payload",
//				Uid:         uuid.New().String(),
//				Subject:     "kor",
//				ErrorMsg:    "error message",
//				Timeout:     100,
//				Retention:   10,
//				CompletedAt: 500,
//				Result:      "empty result",
//				Status:      "running",
//				WorkerUid:   modelWorker.Uid,
//			}
//		_, err = db.NewInsert().Model(&taskMessage).Exec(ctx)
//		req.NoError(err)
//
//		err = taskService.Delete(ctx, taskMessage.Uid)
//		req.NoError(err)
//
//		var task models.Task
//		cnt, err := db.NewSelect().Model(&task).Where("uid = ?", taskMessage.Uid).Count(ctx)
//		req.NoError(err)
//		req.Equal(cnt, 0)
//	})
//}
