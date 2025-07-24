package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetAllLocalBusinesses handles the request to get all local businesses with optional filtering
func (h *LocalHandler) GetAllLocalBusinesses(ctx *fiber.Ctx) error {
	request := local.QueryParamRequestGetLocals{
		City: ctx.Query("city", ""),
		Type: ctx.Query("type", "business"), // Default to business type
	}

	response, err := h.service.GetAllLocalBusinessesWithFilters(ctx.Context(), request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve local businesses",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get all local businesses successful",
		"payload": response,
	})
}

// GetLocalBusinessByID handles the request to get a specific local business by its ID
func (h *LocalHandler) GetLocalBusinessByID(ctx *fiber.Ctx) error {
	businessIDStr := ctx.Params("localBusinessID", "")
	if businessIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing local business ID",
			"message": "Local business ID is required",
		})
	}

	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid local business ID format: %s", businessIDStr),
		})
	}

	response, err := h.service.GetLocalBusinessByID(ctx.Context(), businessID)
	if err != nil {
		if err == local.ErrLBNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Local business not found",
				"message": "The requested local business does not exist",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve local business",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get specific local business successful",
		"payload": response,
	})
}

// CreateLocalBusiness handles the request to create a new local business
func (h *LocalHandler) CreateLocalBusiness(ctx *fiber.Ctx) error {
	var request local.RequestCreateLocalBusiness
	
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

	response, err := h.service.CreateLocalBusiness(ctx.Context(), request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create local business",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "create local business successful",
		"payload": response,
	})
}

// UpdateLocalBusiness handles the request to update an existing local business
func (h *LocalHandler) UpdateLocalBusiness(ctx *fiber.Ctx) error {
	businessIDStr := ctx.Params("localBusinessID", "")
	if businessIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing local business ID",
			"message": "Local business ID is required",
		})
	}

	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid local business ID format: %s", businessIDStr),
		})
	}

	var request local.RequestUpdateLocalBusiness
	
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

	response, err := h.service.UpdateLocalBusiness(ctx.Context(), businessID, request)
	if err != nil {
		if err == local.ErrLBNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Local business not found",
				"message": "The requested local business does not exist",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update local business",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "update local business successful",
		"payload": response,
	})
}

// DeleteLocalBusiness handles the request to delete a local business
func (h *LocalHandler) DeleteLocalBusiness(ctx *fiber.Ctx) error {
	businessIDStr := ctx.Params("localBusinessID", "")
	if businessIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing local business ID",
			"message": "Local business ID is required",
		})
	}

	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid local business ID format: %s", businessIDStr),
		})
	}

	err = h.service.DeleteLocalBusiness(ctx.Context(), businessID)
	if err != nil {
		if err == local.ErrLBNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Local business not found",
				"message": "The requested local business does not exist",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete local business",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete local business successful",
	})
}
