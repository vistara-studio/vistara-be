package rest

import (
	"github.com/vistara-studio/vistara-be/internal/domain/session/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service   service.AuthServiceItf
	validator *validator.Validate
}

func New(service service.AuthServiceItf, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{service: service, validator: validator}
}

func (h *AuthHandler) Mount(router fiber.Router) {
	authGroup := router.Group("/auth")

	authGroup.Post("/register", h.register)
	authGroup.Post("/login", h.login)
}
