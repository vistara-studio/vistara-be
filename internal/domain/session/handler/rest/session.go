package rest

import (
	"time"

	"github.com/vistara-studio/vistara-be/internal/domain/session"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) register(ctx *fiber.Ctx) error {
	var request session.RegisterRequest
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		return err
	}

	response, err := h.service.Register(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *AuthHandler) login(ctx *fiber.Ctx) error {
	var request session.LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		return err
	}

	response, err := h.service.Login(ctx.Context(), request)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login successful",
		"payload": response,
	})
}
