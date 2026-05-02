package httptool

import (
	"github.com/dimazaicev88/ts/internal/dto"

	"github.com/gofiber/fiber/v3"
)

func SendServerErr(ctx fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
		Err: err.Error(),
	})
}

func SendBadRequestErr(ctx fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
		Err: err.Error(),
	})
}
