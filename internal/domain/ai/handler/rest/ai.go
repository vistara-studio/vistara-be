package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vistara-studio/vistara-be/internal/infra/ai"
	"github.com/vistara-studio/vistara-be/internal/middleware"
)

type AIHandler struct {
	aiClient  *ai.Client
	validator *validator.Validate
}

func NewAIHandler(aiClient *ai.Client, validator *validator.Validate) *AIHandler {
	return &AIHandler{
		aiClient:  aiClient,
		validator: validator,
	}
}

func (h *AIHandler) Mount(router fiber.Router) {
	aiGroup := router.Group("/ai")
	aiGroup.Post("/smart-plan", h.GenerateSmartPlan)

	// Service-to-service notification endpoint from vistara-ai
	serviceGroup := router.Group("/service")
	serviceGroup.Use(middleware.ServiceAuthentication())
	serviceGroup.Post("/ai/notify", h.ReceiveNotification)
}

func (h *AIHandler) GenerateSmartPlan(c *fiber.Ctx) error {
	var req ai.SmartPlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	response, err := h.aiClient.GenerateSmartPlan(&req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"message": "AI service is currently unavailable",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Smart plan generated successfully",
		"data":    response,
	})
}

// ReceiveNotification handles notifications from vistara-ai service
func (h *AIHandler) ReceiveNotification(c *fiber.Ctx) error {
	var notification struct {
		Event     string      `json:"event"`
		UserID    *string     `json:"user_id,omitempty"`
		Data      interface{} `json:"data,omitempty"`
		Timestamp string      `json:"timestamp"`
	}

	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid notification body",
		})
	}

	// Process the notification (you can add your business logic here)
	// For example: log the event, update user data, send push notifications, etc.

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification received successfully",
	})
}
