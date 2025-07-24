package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetAllLocalBusinessesWithFilters retrieves all local businesses with optional filtering
func (s *localService) GetAllLocalBusinessesWithFilters(ctx context.Context, request local.QueryParamRequestGetLocals) ([]local.ResponseGetLocalBusinesses, error) {
	localRepository, err := s.repository.NewClient(false)
	if err != nil {
		return []local.ResponseGetLocalBusinesses{}, err
	}

	var localBusinesses []local.Locals
	err = localRepository.GetAllLocalBusinesses(ctx, request, &localBusinesses)
	if err != nil {
		return []local.ResponseGetLocalBusinesses{}, err
	}

	// Transform domain entities to response DTOs
	response := make([]local.ResponseGetLocalBusinesses, len(localBusinesses))
	for i, business := range localBusinesses {
		response[i] = local.ResponseGetLocalBusinesses{
			ID:          business.ID,
			Name:        business.Name,
			Description: business.Description,
			Address:     business.Address,
			City:        business.City,
			Province:    business.Province,
			Longitude:   business.Longitude,
			Latitude:    business.Latitude,
			Label:       business.Label,
			OpenedTime:  business.OpenedTime,
			PhotoUrl:    business.PhotoUrl,
			IsBusiness:  business.IsBusiness,
			CreatedAt:   business.CreatedAt,
		}
	}

	return response, nil
}

// GetLocalBusinessByID retrieves a specific local business by its ID with reviews
func (s *localService) GetLocalBusinessByID(ctx context.Context, businessID uuid.UUID) (local.ResponseGetLocalBusinesses, error) {
	localRepository, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	// Get the local business
	business := &local.Locals{
		ID: businessID,
	}
	err = localRepository.GetLocalBusinessByID(ctx, business)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	// Get reviews for the business
	var reviews []local.Review
	err = localRepository.GetReviewsByLocalBusinessID(ctx, businessID.String(), &reviews)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	// Transform reviews to response format
	reviewResponses := make([]local.ResponseReviews, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = local.ResponseReviews{
			ID:        review.ID,
			Star:      review.Star,
			Content:   review.Content,
			CreatedAt: review.CreatedAt,
			PhotoURL:  review.PhotoURL,
		}
	}

	return local.ResponseGetLocalBusinesses{
		ID:          business.ID,
		Name:        business.Name,
		Description: business.Description,
		Address:     business.Address,
		City:        business.City,
		Province:    business.Province,
		Longitude:   business.Longitude,
		Latitude:    business.Latitude,
		Label:       business.Label,
		OpenedTime:  business.OpenedTime,
		PhotoUrl:    business.PhotoUrl,
		IsBusiness:  business.IsBusiness,
		CreatedAt:   business.CreatedAt,
		Reviews:     reviewResponses,
	}, nil
}

// CreateLocalBusiness creates a new local business
func (s *localService) CreateLocalBusiness(ctx context.Context, request local.RequestCreateLocalBusiness) (local.ResponseGetLocalBusinesses, error) {
	client, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	businessID := uuid.New()
	now := time.Now()

	business := &local.Locals{
		ID:          businessID,
		Name:        request.Name,
		Description: request.Description,
		Address:     request.Address,
		City:        request.City,
		Province:    request.Province,
		Longitude:   request.Longitude,
		Latitude:    request.Latitude,
		Label:       request.Label,
		OpenedTime:  request.OpenedTime,
		PhotoUrl:    request.PhotoUrl,
		IsBusiness:  request.IsBusiness,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = client.CreateLocalBusiness(ctx, business)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	return local.ResponseGetLocalBusinesses{
		ID:          business.ID,
		Name:        business.Name,
		Description: business.Description,
		Address:     business.Address,
		City:        business.City,
		Province:    business.Province,
		Longitude:   business.Longitude,
		Latitude:    business.Latitude,
		Label:       business.Label,
		OpenedTime:  business.OpenedTime,
		PhotoUrl:    business.PhotoUrl,
		IsBusiness:  business.IsBusiness,
		CreatedAt:   business.CreatedAt,
		Reviews:     []local.ResponseReviews{},
	}, nil
}

// UpdateLocalBusiness updates an existing local business
func (s *localService) UpdateLocalBusiness(ctx context.Context, businessID uuid.UUID, request local.RequestUpdateLocalBusiness) (local.ResponseGetLocalBusinesses, error) {
	client, err := s.repository.NewClient(false)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	// First, get the existing business
	business := &local.Locals{ID: businessID}
	err = client.GetLocalBusinessByID(ctx, business)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	// Update only provided fields
	if request.Name != nil {
		business.Name = *request.Name
	}
	if request.Description != nil {
		business.Description = *request.Description
	}
	if request.Address != nil {
		business.Address = *request.Address
	}
	if request.City != nil {
		business.City = *request.City
	}
	if request.Province != nil {
		business.Province = *request.Province
	}
	if request.Longitude != nil {
		business.Longitude = *request.Longitude
	}
	if request.Latitude != nil {
		business.Latitude = *request.Latitude
	}
	if request.Label != nil {
		business.Label = *request.Label
	}
	if request.OpenedTime != nil {
		business.OpenedTime = *request.OpenedTime
	}
	if request.PhotoUrl != nil {
		business.PhotoUrl = *request.PhotoUrl
	}
	if request.IsBusiness != nil {
		business.IsBusiness = *request.IsBusiness
	}

	business.UpdatedAt = time.Now()

	err = client.UpdateLocalBusiness(ctx, business)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	// Get reviews for the response
	var reviews []local.Review
	err = client.GetReviewsByLocalBusinessID(ctx, businessID.String(), &reviews)
	if err != nil {
		return local.ResponseGetLocalBusinesses{}, err
	}

	reviewResponses := make([]local.ResponseReviews, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = local.ResponseReviews{
			ID:        review.ID,
			Star:      review.Star,
			Content:   review.Content,
			CreatedAt: review.CreatedAt,
			PhotoURL:  review.PhotoURL,
		}
	}

	return local.ResponseGetLocalBusinesses{
		ID:          business.ID,
		Name:        business.Name,
		Description: business.Description,
		Address:     business.Address,
		City:        business.City,
		Province:    business.Province,
		Longitude:   business.Longitude,
		Latitude:    business.Latitude,
		Label:       business.Label,
		OpenedTime:  business.OpenedTime,
		PhotoUrl:    business.PhotoUrl,
		IsBusiness:  business.IsBusiness,
		CreatedAt:   business.CreatedAt,
		Reviews:     reviewResponses,
	}, nil
}

// DeleteLocalBusiness deletes a local business
func (s *localService) DeleteLocalBusiness(ctx context.Context, businessID uuid.UUID) error {
	client, err := s.repository.NewClient(false)
	if err != nil {
		return err
	}

	return client.DeleteLocalBusiness(ctx, businessID.String())
}
