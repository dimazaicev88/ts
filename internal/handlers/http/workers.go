package http

import (
	"context"
	"strconv"

	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/httptool"
	"github.com/dimazaicev88/ts/internal/services"
	"github.com/dimazaicev88/ts/internal/storage/filters"
	"github.com/dimazaicev88/ts/internal/tools"

	"github.com/gofiber/fiber/v3"
	"github.com/guregu/null/v6"
)

type Workers struct {
	ctx           context.Context
	fb            *fiber.App
	workerService services.IWorker
}

func NewWorkers(
	ctx context.Context,
	fb *fiber.App,
	workerService services.IWorker,
) Workers {
	return Workers{
		ctx:           ctx,
		fb:            fb,
		workerService: workerService,
	}
}

// FindByName Поиск worker по названию
//
//	@Summary		Поиск worker по названию
//	@Description	Поиск worker по названию
//	@Tags			Workers
//	@Produce		json
//	@Param			workerName	path		string				true	"Название worker"
//	@Success		200			{object}	dto.Worker			"Worker"
//	@Failure		500			{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/wrokers/{workerName} [get]
func (w Workers) FindByName(ctx fiber.Ctx) error {
	result, err := w.workerService.FindByName(w.ctx, ctx.Params("workerName"))
	if err != nil {
		return httptool.SendServerErr(ctx, err)
	}

	return ctx.JSON(result)
}

// Find Поиск worker по названию
//
//	@Summary		Поиск workers
//	@Description	Поиск workers
//	@Tags			Workers
//	@Produce		json
//	@Param			name	query		string				false	"Название worker"
//	@Param			status	query		integer				false	"Статус worker. 0 - offline, 1 - online"
//	@Param			limit	query		integer				false	"Лимит записей. Максимальное кол-во записей 1000"
//	@Param			skip	query		integer				false	"Смещение от первой записи"
//	@Param			count	query		integer				false	"Возвращать кол-во записей"
//	@Success		200		{object}	[]dto.Worker		"Workers"
//	@Failure		500		{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/wrokers [get]
func (w Workers) Find(ctx fiber.Ctx) error {
	var total int
	var err error
	status := ctx.Query("status", "")
	limit, _ := strconv.Atoi(ctx.Query("limit", "1000"))
	if limit > 1000 {
		limit = 1000
	}

	skip, _ := strconv.Atoi(ctx.Query("skip", "0"))
	returnTotal, _ := strconv.Atoi(ctx.Query("count", "0"))

	filterWorker := filters.Worker{
		Skip:  skip,
		Limit: limit,
		Name:  tools.AsNullValue(ctx.Query("name", "")),
	}

	if status == "" {
		filterWorker.Status = tools.AsNullValue(uint8(0))
	} else {
		value, _ := strconv.Atoi(status)
		filterWorker.Status = null.NewValue(uint8(value), true)
	}

	if returnTotal == 1 {
		total, err = w.workerService.Count(w.ctx, filterWorker)
		if err != nil {
			return httptool.SendServerErr(ctx, err)
		}
	}

	workers, err := w.workerService.Find(w.ctx, filterWorker)
	if err != nil {
		return httptool.SendServerErr(ctx, err)
	}

	return ctx.JSON(dto.ResponseWorkers{
		Items: workers,
		Total: total,
	})
}
