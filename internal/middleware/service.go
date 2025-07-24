package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ServiceAuthentication middleware validates service-to-service requests
func ServiceAuthentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for service header
		serviceHeader := c.Get("X-Service")
		if serviceHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Service authentication required",
			})
		}

		// Validate that the service is vistara-ai
		if serviceHeader != "vistara-ai" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized service",
			})
		}

		return c.Next()
	}
}
