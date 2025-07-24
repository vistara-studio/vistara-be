package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ServiceAuthentication middleware validates service-to-service requests
// Only allows requests from vistara-ai service with proper X-Service header
func ServiceAuthentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for service authentication header
		serviceHeader := c.Get("X-Service")
		if serviceHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Service authentication required",
			})
		}

		// Validate that the request is from vistara-ai
		if serviceHeader != "vistara-ai" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized service",
			})
		}

		return c.Next()
	}
}
