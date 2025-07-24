package http

import (
	"time"

	"github.com/vistara-studio/vistara-be/pkg/cerr"
	_validator "github.com/vistara-studio/vistara-be/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ErrorHandler: ErrorHandler(),
	})
}

func ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		if ce, ok := err.(*cerr.CustomError); ok {
			return ctx.Status(ce.Code).JSON(fiber.Map{
				"message": ce.Message,
				"error":   ce.Error(),
			})
		}

		if ve, ok := err.(validator.ValidationErrors); ok {
			errorList := _validator.GetError(err, ve)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "validation error",
				"error":   errorList,
			})
		}

		if fe, ok := err.(*fiber.Error); ok {
			return ctx.Status(fe.Code).JSON(fiber.Map{
				"message": fe.Message,
				"error":   fe.Message,
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fiber.StatusInternalServerError,
			"error":   err.Error(),
		})

	}
}
