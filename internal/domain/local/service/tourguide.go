package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetAllTouristAttractions retrieves all tourist attractions with optional city filtering
func (s *localService) GetAllTouristAttractions(ctx context.Context, city string) ([]local.ResponseGetTourGuide, error) {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return []local.ResponseGetTourGuide{}, err
	}

	var touristAttractions []local.TouristAttractions
	err = repository.GetAllTouristAttractions(ctx, city, &touristAttractions)
	if err != nil {
		return []local.ResponseGetTourGuide{}, err
	}

	// Transform domain entities to response DTOs
	response := make([]local.ResponseGetTourGuide, len(touristAttractions))
	for i, attraction := range touristAttractions {
		response[i] = local.ResponseGetTourGuide{
			ID:                          attraction.ID,
			Name:                        attraction.Name,
			Description:                 attraction.Description,
			Address:                     attraction.Address,
			City:                        attraction.City,
			Province:                    attraction.Province,
			Longitude:                   attraction.Longitude,
			Latitude:                    attraction.Latitude,
			PhotoUrl:                    attraction.PhotoURL,
			TourGuidePrice:              attraction.TourGuidePrice,
			TourGuideCount:              attraction.TourGuideCount,
			TourGuideDiscountPercentage: attraction.TourGuideDiscountPercentage,
			Price:                       attraction.Price,
			DiscountPercentage:          attraction.DiscountPercentage,
			CreatedAt:                   attraction.CreatedAt,
		}
	}

	return response, nil
}

// GetTouristAttractionByID retrieves a specific tourist attraction by its ID with bookings/reviews
func (s *localService) GetTouristAttractionByID(ctx context.Context, attractionID uuid.UUID) (local.ResponseGetTourGuide, error) {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGetTourGuide{}, err
	}

	// Get the tourist attraction
	attraction := &local.TouristAttractions{
		ID: attractionID,
	}
	err = repository.GetTouristAttractionByID(ctx, attraction)
	if err != nil {
		return local.ResponseGetTourGuide{}, err
	}

	// Get bookings for the attraction (used as reviews)
	var bookings []local.TourGuideBookings
	err = repository.GetBookingsByTouristAttractionID(ctx, attractionID.String(), &bookings)
	if err != nil {
		return local.ResponseGetTourGuide{}, err
	}

	// Transform bookings to review format
	reviews := make([]local.ResponseReviews, len(bookings))
	for i, booking := range bookings {
		reviews[i] = local.ResponseReviews{
			ID:        booking.ID,
			Star:      booking.Star,
			Content:   booking.Content,
			CreatedAt: booking.CreatedAt,
			PhotoURL:  booking.PhotoURL,
		}
	}

	return local.ResponseGetTourGuide{
		ID:                          attraction.ID,
		Name:                        attraction.Name,
		Description:                 attraction.Description,
		Address:                     attraction.Address,
		City:                        attraction.City,
		Province:                    attraction.Province,
		Longitude:                   attraction.Longitude,
		Latitude:                    attraction.Latitude,
		PhotoUrl:                    attraction.PhotoURL,
		TourGuidePrice:              attraction.TourGuidePrice,
		TourGuideCount:              attraction.TourGuideCount,
		TourGuideDiscountPercentage: attraction.TourGuideDiscountPercentage,
		Price:                       attraction.Price,
		DiscountPercentage:          attraction.DiscountPercentage,
		CreatedAt:                   attraction.CreatedAt,
		Reviews:                     reviews,
	}, nil
}

// GetFullyBookedDates retrieves dates when the tourist attraction is fully booked
func (s *localService) GetFullyBookedDates(ctx context.Context, attractionID string, year, month int) ([]string, error) {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return []string{}, err
	}

	var dates []string
	err = repository.GetFullyBookedDates(ctx, attractionID, year, month, &dates)
	if err != nil {
		return []string{}, err
	}

	return dates, nil
}

// GeneratePaymentSnapLink generates a Midtrans Snap payment link for tour guide booking
func (s *localService) GeneratePaymentSnapLink(ctx context.Context, request local.RequestGenerateSnapLink) (local.ResponseGenerateSnapLink, error) {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, err
	}

	// Parse and validate tourist attraction ID
	attractionID, err := uuid.Parse(request.TAID)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, fmt.Errorf("invalid tourist attraction ID: %w", err)
	}

	// Get tourist attraction details
	attraction := &local.TouristAttractions{
		ID: attractionID,
	}
	err = repository.GetTouristAttractionByID(ctx, attraction)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, err
	}

	// Generate unique transaction ID
	transactionID, err := uuid.NewV7()
	if err != nil {
		return local.ResponseGenerateSnapLink{}, fmt.Errorf("failed to generate transaction ID: %w", err)
	}

	// Create Midtrans Snap request
	snapRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionID.String(),
			GrossAmt: attraction.TourGuidePrice,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Expiry: &snap.ExpiryDetails{
			Duration: 24, // 24 hours expiry
			Unit:     "hours",
		},
	}

	// Create payment token
	snapToken, err := s.snapClient.CreateTransactionToken(snapRequest)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, fmt.Errorf("failed to create payment token: %w", err)
	}

	// Parse booking date
	layout := "2006-01-02"
	bookedAt, err := time.Parse(layout, request.BookedAt)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, fmt.Errorf("invalid booking date format: %w", err)
	}
	bookedAt = time.Date(bookedAt.Year(), bookedAt.Month(), bookedAt.Day(), 0, 0, 0, 0, time.UTC)

	// Parse user ID
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create booking record
	booking := &local.TourGuideBookings{
		ID:                   transactionID,
		PaymentURL:           fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v4/redirection/%s", snapToken),
		BookedAt:             bookedAt,
		Status:               "pending_payment",
		UserID:               userID,
		TouristAttractionsID: attractionID,
	}

	err = repository.CreateTourGuideBooking(ctx, booking)
	if err != nil {
		return local.ResponseGenerateSnapLink{}, fmt.Errorf("failed to create booking: %w", err)
	}

	return local.ResponseGenerateSnapLink{
		TAID:       attractionID.String(),
		PaymentUrl: snapToken,
	}, nil
}

// CreateTouristAttraction creates a new tourist attraction
func (s *localService) CreateTouristAttraction(ctx context.Context, request local.RequestCreateTouristAttraction) (local.ResponseGetTourGuide, error) {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGetTourGuide{}, err
	}

	// Generate new UUID for the attraction
	attractionID := uuid.New()
	now := time.Now()

	attraction := &local.TouristAttractions{
		ID:                          attractionID,
		Name:                        request.Name,
		Description:                 request.Description,
		Address:                     request.Address,
		City:                        request.City,
		Province:                    request.Province,
		Longitude:                   request.Longitude,
		Latitude:                    request.Latitude,
		PhotoURL:                    request.PhotoUrl,
		TourGuidePrice:              request.TourGuidePrice,
		TourGuideCount:              request.TourGuideCount,
		TourGuideDiscountPercentage: request.TourGuideDiscountPercentage,
		Price:                       request.Price,
		DiscountPercentage:          request.DiscountPercentage,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}

	err = repository.CreateTouristAttraction(ctx, attraction)
	if err != nil {
		return local.ResponseGetTourGuide{}, fmt.Errorf("failed to create tourist attraction: %w", err)
	}

	// Return the created attraction
	return local.ResponseGetTourGuide{
		ID:                          attraction.ID,
		Name:                        attraction.Name,
		Description:                 attraction.Description,
		Address:                     attraction.Address,
		City:                        attraction.City,
		Province:                    attraction.Province,
		Longitude:                   attraction.Longitude,
		Latitude:                    attraction.Latitude,
		PhotoUrl:                    attraction.PhotoURL,
		TourGuidePrice:              attraction.TourGuidePrice,
		TourGuideCount:              attraction.TourGuideCount,
		TourGuideDiscountPercentage: attraction.TourGuideDiscountPercentage,
		Price:                       attraction.Price,
		DiscountPercentage:          attraction.DiscountPercentage,
		CreatedAt:                   attraction.CreatedAt,
		Reviews:                     []local.ResponseReviews{}, // Empty reviews for new attraction
	}, nil
}

// UpdateTouristAttraction updates an existing tourist attraction
func (s *localService) UpdateTouristAttraction(ctx context.Context, attractionID uuid.UUID, request local.RequestUpdateTouristAttraction) (local.ResponseGetTourGuide, error) {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGetTourGuide{}, err
	}

	// First, get the existing attraction
	attraction := &local.TouristAttractions{ID: attractionID}
	err = repository.GetTouristAttractionByID(ctx, attraction)
	if err != nil {
		return local.ResponseGetTourGuide{}, fmt.Errorf("failed to get tourist attraction: %w", err)
	}

	// Update only provided fields
	if request.Name != nil {
		attraction.Name = *request.Name
	}
	if request.Description != nil {
		attraction.Description = *request.Description
	}
	if request.Address != nil {
		attraction.Address = *request.Address
	}
	if request.City != nil {
		attraction.City = *request.City
	}
	if request.Province != nil {
		attraction.Province = *request.Province
	}
	if request.Longitude != nil {
		attraction.Longitude = *request.Longitude
	}
	if request.Latitude != nil {
		attraction.Latitude = *request.Latitude
	}
	if request.PhotoUrl != nil {
		attraction.PhotoURL = *request.PhotoUrl
	}
	if request.TourGuidePrice != nil {
		attraction.TourGuidePrice = *request.TourGuidePrice
	}
	if request.TourGuideCount != nil {
		attraction.TourGuideCount = *request.TourGuideCount
	}
	if request.TourGuideDiscountPercentage != nil {
		attraction.TourGuideDiscountPercentage = *request.TourGuideDiscountPercentage
	}
	if request.Price != nil {
		attraction.Price = *request.Price
	}
	if request.DiscountPercentage != nil {
		attraction.DiscountPercentage = *request.DiscountPercentage
	}

	attraction.UpdatedAt = time.Now()

	err = repository.UpdateTouristAttraction(ctx, attraction)
	if err != nil {
		return local.ResponseGetTourGuide{}, fmt.Errorf("failed to update tourist attraction: %w", err)
	}

	// Return the updated attraction
	return local.ResponseGetTourGuide{
		ID:                          attraction.ID,
		Name:                        attraction.Name,
		Description:                 attraction.Description,
		Address:                     attraction.Address,
		City:                        attraction.City,
		Province:                    attraction.Province,
		Longitude:                   attraction.Longitude,
		Latitude:                    attraction.Latitude,
		PhotoUrl:                    attraction.PhotoURL,
		TourGuidePrice:              attraction.TourGuidePrice,
		TourGuideCount:              attraction.TourGuideCount,
		TourGuideDiscountPercentage: attraction.TourGuideDiscountPercentage,
		Price:                       attraction.Price,
		DiscountPercentage:          attraction.DiscountPercentage,
		CreatedAt:                   attraction.CreatedAt,
		Reviews:                     []local.ResponseReviews{}, // Reviews would need separate loading
	}, nil
}

// DeleteTouristAttraction deletes a tourist attraction
func (s *localService) DeleteTouristAttraction(ctx context.Context, attractionID uuid.UUID) error {
	repository, err := s.repository.NewClient(false)
	if err != nil {
		return err
	}

	err = repository.DeleteTouristAttraction(ctx, attractionID.String())
	if err != nil {
		return fmt.Errorf("failed to delete tourist attraction: %w", err)
	}

	return nil
}
