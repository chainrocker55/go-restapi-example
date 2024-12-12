package errors

import (
	"creditlimit-connector/app/log"
	"creditlimit-connector/app/models"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	code := "500"
	message := fiber.ErrInternalServerError.Message

	log.Error(err)

	if e, ok := err.(*fiber.Error); ok {
		message = e.Message
	} else if e, ok := err.(*url.Error); ok {
		statusCode = fiber.StatusBadGateway
		message = e.Err.Error()
	} else if e, ok := err.(*models.ErrorResponse); ok {
		statusCode = e.HttpStatus
		code = e.Code
		message = e.Message
	}

	ctx.Status(statusCode).JSON(fiber.Map{"code": code, "message": message})

	return nil
}
