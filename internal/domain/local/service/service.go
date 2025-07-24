package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
	"github.com/vistara-studio/vistara-be/internal/domain/local/repository"
)

// localService implements the local business service
type localService struct {
	repository repository.RepositoryInterface
	snapClient snap.Client
	coreAPI    coreapi.Client
}

// LocalServiceInterface defines the contract for local business operations
type LocalServiceInterface interface {
	// Local business operations
	GetAllLocalBusinessesWithFilters(ctx context.Context, request local.QueryParamRequestGetLocals) ([]local.ResponseGetLocalBusinesses, error)
	GetLocalBusinessByID(ctx context.Context, businessID uuid.UUID) (local.ResponseGetLocalBusinesses, error)
	CreateLocalBusiness(ctx context.Context, request local.RequestCreateLocalBusiness) (local.ResponseGetLocalBusinesses, error)
	UpdateLocalBusiness(ctx context.Context, businessID uuid.UUID, request local.RequestUpdateLocalBusiness) (local.ResponseGetLocalBusinesses, error)
	DeleteLocalBusiness(ctx context.Context, businessID uuid.UUID) error
	
	// Tourist attraction operations
	GetAllTouristAttractions(ctx context.Context, city string) ([]local.ResponseGetTourGuide, error)
	GetTouristAttractionByID(ctx context.Context, attractionID uuid.UUID) (local.ResponseGetTourGuide, error)
	CreateTouristAttraction(ctx context.Context, request local.RequestCreateTouristAttraction) (local.ResponseGetTourGuide, error)
	UpdateTouristAttraction(ctx context.Context, attractionID uuid.UUID, request local.RequestUpdateTouristAttraction) (local.ResponseGetTourGuide, error)
	DeleteTouristAttraction(ctx context.Context, attractionID uuid.UUID) error
	
	// Booking operations
	GeneratePaymentSnapLink(ctx context.Context, request local.RequestGenerateSnapLink) (local.ResponseGenerateSnapLink, error)
	GetFullyBookedDates(ctx context.Context, attractionID string, year, month int) ([]string, error)
}

// New creates a new local service instance
func New(repo repository.RepositoryInterface, snapClient snap.Client, coreAPI coreapi.Client) LocalServiceInterface {
	return &localService{
		repository: repo,
		snapClient: snapClient,
		coreAPI:    coreAPI,
	}
}
