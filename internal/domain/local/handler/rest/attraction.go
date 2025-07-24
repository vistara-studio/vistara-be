package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// CreateTouristAttraction handles the request to create a new tourist attraction
func (h *LocalHandler) CreateTouristAttraction(ctx *fiber.Ctx) error {
	var request local.RequestCreateTouristAttraction
	
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse JSON request body",
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "validation error",
		})
	}

	response, err := h.service.CreateTouristAttraction(ctx.Context(), request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create tourist attraction",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "create tourist attraction successful",
		"payload": response,
	})
}

// UpdateTouristAttraction handles the request to update an existing tourist attraction
func (h *LocalHandler) UpdateTouristAttraction(ctx *fiber.Ctx) error {
	attractionIDStr := ctx.Params("attractionID", "")
	if attractionIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing tourist attraction ID",
			"message": "Tourist attraction ID is required",
		})
	}

	attractionID, err := uuid.Parse(attractionIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid tourist attraction ID format: %s", attractionIDStr),
		})
	}

	var request local.RequestUpdateTouristAttraction
	
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse JSON request body",
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "validation error",
		})
	}

	response, err := h.service.UpdateTouristAttraction(ctx.Context(), attractionID, request)
	if err != nil {
		if err == local.ErrLBNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Tourist attraction not found",
				"message": "The requested tourist attraction does not exist",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update tourist attraction",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "update tourist attraction successful",
		"payload": response,
	})
}

// DeleteTouristAttraction handles the request to delete a tourist attraction
func (h *LocalHandler) DeleteTouristAttraction(ctx *fiber.Ctx) error {
	attractionIDStr := ctx.Params("attractionID", "")
	if attractionIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing tourist attraction ID",
			"message": "Tourist attraction ID is required",
		})
	}

	attractionID, err := uuid.Parse(attractionIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid tourist attraction ID format: %s", attractionIDStr),
		})
	}

	err = h.service.DeleteTouristAttraction(ctx.Context(), attractionID)
	if err != nil {
		if err == local.ErrLBNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Tourist attraction not found",
				"message": "The requested tourist attraction does not exist",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete tourist attraction",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete tourist attraction successful",
	})
}
