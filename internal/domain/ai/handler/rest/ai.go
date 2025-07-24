package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vistara-studio/vistara-be/internal/infra/ai"
	"github.com/vistara-studio/vistara-be/internal/middleware"
	"github.com/vistara-studio/vistara-be/pkg/jwt"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiClient  *ai.Client
	validator *validator.Validate
	jwt       *jwt.JWTStruct
}

// NewAIHandler creates a new AI handler instance
func NewAIHandler(aiClient *ai.Client, validator *validator.Validate, jwt *jwt.JWTStruct) *AIHandler {
	return &AIHandler{
		aiClient:  aiClient,
		validator: validator,
		jwt:       jwt,
	}
}

// Mount registers AI routes on the router
func (h *AIHandler) Mount(router fiber.Router) {
	// Protected AI endpoints (requires JWT authentication)
	aiGroup := router.Group("/ai")
	aiGroup.Use(middleware.Authentication(h.jwt))
	aiGroup.Post("/smart-planner", h.GenerateSmartPlan)

	// Service-to-service endpoints (requires service authentication)  
	serviceGroup := router.Group("/service")
	serviceGroup.Use(middleware.ServiceAuthentication())
	serviceGroup.Post("/ai/notify", h.ReceiveNotification)
}

// GenerateSmartPlan handles smart travel plan generation requests
func (h *AIHandler) GenerateSmartPlan(c *fiber.Ctx) error {
	// Get user information from JWT token
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user context",
		})
	}
	
	isPremiumInterface := c.Locals("is_premium")
	isPremium, ok := isPremiumInterface.(bool)
	if !ok {
		// Default to false if not present or invalid
		isPremium = false
	}
	
	var req ai.SmartPlanRequest
	
	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// Add user context to the request
	req.UserID = &userID
	req.IsPremium = &isPremium

	// Call AI service
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
		"user_id": userID,
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

	// Parse notification body
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid notification body",
		})
	}

	// Process notification (extend with business logic as needed)
	// Examples: log events, update user data, send push notifications

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification received successfully",
	})
}
