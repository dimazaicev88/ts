package http

import (
	"context"
	"strconv"

	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/httptool"
	"github.com/dimazaicev88/ts/internal/services"
	"github.com/dimazaicev88/ts/internal/storage/filters"

	"github.com/gofiber/fiber/v3"
	"github.com/segmentio/encoding/json"
)

type Tasks struct {
	ctx         context.Context
	fb          *fiber.App
	taskService services.ITask
}

func NewTasks(
	ctx context.Context,
	fb *fiber.App,
	taskService services.ITask,
) Tasks {
	return Tasks{
		ctx:         ctx,
		fb:          fb,
		taskService: taskService,
	}
}

// FindTask Поиск задач
//
//	@Summary		Поиск задач
//	@Description	Поиск задач
//	@Tags			Tasks
//	@Produce		json
//	@Param			workerUID	path		string				true	"Uid worker"
//	@Param			status		query		string				false	"Статус задачи"
//	@Param			limit		query		integer				false	"Лимит записей. Максимальное кол-во записей 1000"
//	@Param			skip		query		integer				false	"Смещение от первой записи"
//	@Param			count		query		integer				false	"Флаг указывают нужно ли возвращать общее кол-во записей, 0 - нет, 1 - да"
//	@Success		200			{object}	dto.ResponseTasks	"Список задач"
//	@Failure		500			{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/tasks/workers/{workerUID} [get]
func (t Tasks) FindTask(ctx fiber.Ctx) error {
	var (
		err    error
		result []dto.ResponseTask
		total  int
	)

	limit, _ := strconv.Atoi(ctx.Query("limit", "1000"))
	if limit > 1000 {
		limit = 1000
	}

	skip, _ := strconv.Atoi(ctx.Query("skip", "0"))
	returnTotal, _ := strconv.Atoi(ctx.Query("count", "0"))
	filter := filters.Tasks{
		Status:    ctx.Query("status", ""),
		Limit:     limit,
		Skip:      skip,
		WorkerUid: ctx.Params("workerUid"),
	}

	if returnTotal == 1 {
		total, err = t.taskService.Count(t.ctx, filter)
		if err != nil {
			return httptool.SendServerErr(ctx, err)
		}
	}

	result, err = t.taskService.Find(t.ctx, filter)
	if err != nil {
		return httptool.SendServerErr(ctx, err)
	}

	return ctx.JSON(dto.ResponseTasks{
		Items: result,
		Total: total,
	})
}

// StopTask Остановить задачу.
//
//	@Summary		Остановить задачу.
//	@Description	Остановить задачу можно только в статусе active.
//	@Description	Отправка команды остановки задачи не гарантируем остановку задачи
//	@Tags			Tasks
//	@Produce		json
//
//	@Param			uid	path	string	true	"UID задачи"
//
//	@Success		204	"Успешная посылка команды остановки задачи"
//	@Failure		500	{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/tasks/{uid}/stop [post]
func (t Tasks) StopTask(ctx fiber.Ctx) {
	//TODO послать в centrifuge команду остановки задачи
}

// DeleteTask Удалить задачу.
//
//	@Summary		Удалить задачи.
//	@Description	Удалить задачи можно только в статусе complete, archived.
//	@Tags			Tasks
//	@Produce		json
//	@Param			uids	body	[]dto.ListUID	true	"Список UID задач"
//	@Success		204		"Успешное удаление задач"
//	@Failure		500		{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/tasks [delete]
func (t Tasks) DeleteTask(ctx fiber.Ctx) error {
	var listUID dto.ListUID
	if err := json.Unmarshal(ctx.Body(), &listUID); err != nil {
		return httptool.SendServerErr(ctx, err)
	}
	err := t.taskService.Delete(t.ctx, listUID.UIDS)
	if err != nil {
		return httptool.SendServerErr(ctx, err)
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
