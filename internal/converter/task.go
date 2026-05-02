package converter

import (
	"github.com/dimazaicev88/ts/internal/dto"
	"github.com/dimazaicev88/ts/internal/storage/models"
)

func TaskModelToDto(model models.Task) dto.ResponseTask {
	return dto.ResponseTask{
		Uid:         model.Uid,
		Payload:     model.Payload,
		Status:      model.Status,
		DateCreate:  model.DateCreate,
		DateUpdate:  model.DateUpdate,
		Result:      model.Result,
		Subject:     model.Subject,
		WorkerUid:   model.WorkerUid,
		Timeout:     model.Timeout,
		Retention:   model.Retention,
		CompletedAt: model.CompletedAt,
		ErrorMsg:    model.ErrorMsg,
	}
}

func ListTaskModelToListDto(models []models.Task) []dto.ResponseTask {
	var listDto []dto.ResponseTask
	for _, model := range models {
		listDto = append(listDto, TaskModelToDto(model))
	}
	return listDto
}
