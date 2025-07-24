package rest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetAllTouristAttractions handles the request to get all tourist attractions
func (h *LocalHandler) GetAllTouristAttractions(ctx *fiber.Ctx) error {
	city := ctx.Query("city", "")

	response, err := h.service.GetAllTouristAttractions(ctx.Context(), city)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve tourist attractions",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get all tour guides successful",
		"payload": response,
	})
}

// GetTouristAttractionByID handles the request to get a specific tourist attraction by ID
func (h *LocalHandler) GetTouristAttractionByID(ctx *fiber.Ctx) error {
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

	response, err := h.service.GetTouristAttractionByID(ctx.Context(), attractionID)
	if err != nil {
		if err == local.ErrLBNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Tourist attraction not found",
				"message": "The requested tourist attraction does not exist",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve tourist attraction",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get specific tourist attraction successful",
		"payload": response,
	})
}

// GetFullyBookedDates handles the request to get fully booked dates for a tourist attraction
func (h *LocalHandler) GetFullyBookedDates(ctx *fiber.Ctx) error {
	attractionIDStr := ctx.Params("attractionID", "")
	if attractionIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing tourist attraction ID",
			"message": "Tourist attraction ID is required",
		})
	}

	_, err := uuid.Parse(attractionIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid tourist attraction ID format: %s", attractionIDStr),
		})
	}

	// Get year and month from query params, default to current year and month
	now := time.Now()
	yearStr := ctx.Query("year", strconv.Itoa(now.Year()))
	monthStr := ctx.Query("month", strconv.Itoa(int(now.Month())))

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid year format",
			"message": "Year must be a valid integer",
		})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid month format",
			"message": "Month must be an integer between 1 and 12",
		})
	}

	response, err := h.service.GetFullyBookedDates(ctx.Context(), attractionIDStr, year, month)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve fully booked dates",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get fully booked dates successful",
		"payload": response,
	})
}

// CreateTourGuideBooking handles the request to create a tour guide booking with payment
func (h *LocalHandler) CreateTourGuideBooking(ctx *fiber.Ctx) error {
	var request local.RequestGenerateSnapLink
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": err.Error(),
		})
	}

	attractionIDStr := ctx.Params("attractionID", "")
	if attractionIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing tourist attraction ID",
			"message": "Tourist attraction ID is required",
		})
	}

	_, err := uuid.Parse(attractionIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"message": fmt.Sprintf("Invalid tourist attraction ID format: %s", attractionIDStr),
		})
	}

	// Get user ID from JWT token
	userIDRaw, ok := ctx.Locals("user_id").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authentication required",
			"message": "Failed to get user ID from authentication token",
		})
	}

	_, err = uuid.Parse(userIDRaw)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Invalid user ID format in authentication token",
		})
	}

	// Set IDs from path and auth context
	request.TAID = attractionIDStr
	request.UserID = userIDRaw

	response, err := h.service.GeneratePaymentSnapLink(ctx.Context(), request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate payment link",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "payment link generated successfully",
		"payload": response,
	})
}
