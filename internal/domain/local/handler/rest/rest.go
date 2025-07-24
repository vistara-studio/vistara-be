package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vistara-studio/vistara-be/internal/domain/local/service"
	"github.com/vistara-studio/vistara-be/internal/middleware"
	"github.com/vistara-studio/vistara-be/pkg/jwt"
)

// LocalHandler handles HTTP requests for local business and tourist attraction endpoints
type LocalHandler struct {
	service   service.LocalServiceInterface
	validator *validator.Validate
	jwt       *jwt.JWTStruct
}

// New creates a new LocalHandler instance
func New(service service.LocalServiceInterface, validator *validator.Validate, jwt *jwt.JWTStruct) *LocalHandler {
	return &LocalHandler{
		service:   service,
		validator: validator,
		jwt:       jwt,
	}
}

// Mount registers all local business and tourist attraction routes
func (h *LocalHandler) Mount(router fiber.Router) {
	// Service-to-service routes for AI integration (with service authentication)
	serviceGroup := router.Group("/service")
	serviceGroup.Use(middleware.ServiceAuthentication())
	serviceGroup.Get("/locals", h.GetAllLocalBusinesses)
	serviceGroup.Get("/tourist-attractions", h.GetAllTouristAttractions)

	// Local business routes - All require authentication
	localGroup := router.Group("/locals")
	localGroup.Use(middleware.Authentication(h.jwt))
	localGroup.Get("/", h.GetAllLocalBusinesses)
	localGroup.Get("/:localBusinessID", h.GetLocalBusinessByID)
	localGroup.Post("/", h.CreateLocalBusiness)
	localGroup.Put("/:localBusinessID", h.UpdateLocalBusiness)
	localGroup.Delete("/:localBusinessID", h.DeleteLocalBusiness)

	// Tourist attraction routes - All require authentication
	attractionGroup := router.Group("/tourist-attractions")
	attractionGroup.Use(middleware.Authentication(h.jwt))
	attractionGroup.Get("/", h.GetAllTouristAttractions)
	attractionGroup.Get("/:attractionID", h.GetTouristAttractionByID)
	attractionGroup.Post("/", h.CreateTouristAttraction)
	attractionGroup.Put("/:attractionID", h.UpdateTouristAttraction)
	attractionGroup.Delete("/:attractionID", h.DeleteTouristAttraction)
	attractionGroup.Get("/:attractionID/availability", h.GetFullyBookedDates)
	attractionGroup.Post("/:attractionID/book", h.CreateTourGuideBooking)
}
