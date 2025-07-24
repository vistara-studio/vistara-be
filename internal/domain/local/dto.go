package local

import (
	"time"

	"github.com/google/uuid"
)

type ResponseGetLocalBusinesses struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Address     string            `json:"address"`
	City        string            `json:"city"`
	Province    string            `json:"province"`
	Longitude   string            `json:"longitude"`
	Latitude    string            `json:"latitude"`
	Label       string            `json:"label"`
	OpenedTime  string            `json:"opened_time"`
	PhotoUrl    string            `json:"photo_url"`
	IsBusiness  bool              `json:"is_business"`
	CreatedAt   time.Time         `json:"created_at"`
	Reviews     []ResponseReviews `json:"reviews,omitempty"`
}

type ResponseGetTourGuide struct {
	ID                          uuid.UUID         `json:"id"`
	Name                        string            `json:"name"`
	Description                 string            `json:"description"`
	Address                     string            `json:"address"`
	City                        string            `json:"city"`
	Province                    string            `json:"province"`
	Longitude                   float64           `json:"longitude"`
	Latitude                    float64           `json:"latitude"`
	PhotoUrl                    string            `json:"photo_url"`
	TourGuidePrice              int64             `json:"tour_guide_price"`
	TourGuideCount              int               `json:"tour_guide_count"`
	TourGuideDiscountPercentage float32           `json:"tour_guide_discount_percentage"`
	Price                       int64             `json:"price"`
	DiscountPercentage          float32           `json:"discount_percentage"`
	CreatedAt                   time.Time         `json:"created_at"`
	Reviews                     []ResponseReviews `json:"reviews,omitempty"`
}

type ResponseReviews struct {
	ID        uuid.UUID `json:"id"`
	Star      int       `json:"star"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	PhotoURL  string    `json:"photo_url"`
}

type QueryParamRequestGetLocals struct {
	City string
	Type string
}

// RequestCreateLocalBusiness represents the request body for creating a new local business
type RequestCreateLocalBusiness struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	Address     string  `json:"address" validate:"required,min=10,max=200"`
	City        string  `json:"city" validate:"required,min=2,max=50"`
	Province    string  `json:"province" validate:"required,min=2,max=50"`
	Longitude   string  `json:"longitude" validate:"required"`
	Latitude    string  `json:"latitude" validate:"required"`
	Label       string  `json:"label" validate:"required,min=2,max=50"`
	OpenedTime  string  `json:"opened_time" validate:"required"`
	PhotoUrl    string  `json:"photo_url" validate:"required,url"`
	IsBusiness  bool    `json:"is_business"`
}

// RequestUpdateLocalBusiness represents the request body for updating a local business
type RequestUpdateLocalBusiness struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=10,max=500"`
	Address     *string `json:"address,omitempty" validate:"omitempty,min=10,max=200"`
	City        *string `json:"city,omitempty" validate:"omitempty,min=2,max=50"`
	Province    *string `json:"province,omitempty" validate:"omitempty,min=2,max=50"`
	Longitude   *string `json:"longitude,omitempty"`
	Latitude    *string `json:"latitude,omitempty"`
	Label       *string `json:"label,omitempty" validate:"omitempty,min=2,max=50"`
	OpenedTime  *string `json:"opened_time,omitempty"`
	PhotoUrl    *string `json:"photo_url,omitempty" validate:"omitempty,url"`
	IsBusiness  *bool   `json:"is_business,omitempty"`
}

// RequestCreateTouristAttraction represents the request body for creating a new tourist attraction
type RequestCreateTouristAttraction struct {
	Name                         string  `json:"name" validate:"required,min=3,max=100"`
	Description                  string  `json:"description" validate:"required,min=10,max=500"`
	Address                      string  `json:"address" validate:"required,min=10,max=200"`
	City                         string  `json:"city" validate:"required,min=2,max=50"`
	Province                     string  `json:"province" validate:"required,min=2,max=50"`
	Longitude                    float64 `json:"longitude" validate:"required"`
	Latitude                     float64 `json:"latitude" validate:"required"`
	PhotoUrl                     string  `json:"photo_url" validate:"required,url"`
	TourGuidePrice               int64   `json:"tour_guide_price" validate:"required,min=0"`
	TourGuideCount               int     `json:"tour_guide_count" validate:"required,min=1"`
	TourGuideDiscountPercentage  float32 `json:"tour_guide_discount_percentage" validate:"min=0,max=100"`
	Price                        int64   `json:"price" validate:"required,min=0"`
	DiscountPercentage           float32 `json:"discount_percentage" validate:"min=0,max=100"`
}

// RequestUpdateTouristAttraction represents the request body for updating a tourist attraction
type RequestUpdateTouristAttraction struct {
	Name                         *string  `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description                  *string  `json:"description,omitempty" validate:"omitempty,min=10,max=500"`
	Address                      *string  `json:"address,omitempty" validate:"omitempty,min=10,max=200"`
	City                         *string  `json:"city,omitempty" validate:"omitempty,min=2,max=50"`
	Province                     *string  `json:"province,omitempty" validate:"omitempty,min=2,max=50"`
	Longitude                    *float64 `json:"longitude,omitempty"`
	Latitude                     *float64 `json:"latitude,omitempty"`
	PhotoUrl                     *string  `json:"photo_url,omitempty" validate:"omitempty,url"`
	TourGuidePrice               *int64   `json:"tour_guide_price,omitempty" validate:"omitempty,min=0"`
	TourGuideCount               *int     `json:"tour_guide_count,omitempty" validate:"omitempty,min=1"`
	TourGuideDiscountPercentage  *float32 `json:"tour_guide_discount_percentage,omitempty" validate:"omitempty,min=0,max=100"`
	Price                        *int64   `json:"price,omitempty" validate:"omitempty,min=0"`
	DiscountPercentage           *float32 `json:"discount_percentage,omitempty" validate:"omitempty,min=0,max=100"`
}

type RequestGenerateSnapLink struct {
	UserID   string ``
	TAID     string ``
	BookedAt string `json:"booked_at" validate:"required"`
}

type ResponseGenerateSnapLink struct {
	TAID       string `json:"ta_id"`
	PaymentUrl string `json:"payment_url"`
}
