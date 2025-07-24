package local

import (
	"time"

	"github.com/google/uuid"
)

type TourGuideBookings struct {
	ID                   uuid.UUID `db:"id"`
	PaymentURL           string    `db:"payment_url"`
	Star                 int       `db:"star"`
	Content              string    `db:"content"`
	BookedAt             time.Time `db:"booked_at"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
	Status               string    `db:"status"`
	UserID               uuid.UUID `db:"user_id"`
	TouristAttractionsID uuid.UUID `db:"tourist_attraction_id"`
	PhotoURL             string    `db:"photo_url"`
}

type TouristAttractions struct {
	ID                          uuid.UUID `db:"id"`
	Name                        string    `db:"name"`
	Description                 string    `db:"description"`
	Address                     string    `db:"address"`
	City                        string    `db:"city"`
	Province                    string    `db:"province"`
	Longitude                   float64   `db:"longitude"`
	Latitude                    float64   `db:"latitude"`
	PhotoURL                    string    `db:"photo_url"`
	TourGuidePrice              int64     `db:"tour_guide_price"`
	TourGuideCount              int       `db:"tour_guide_count"`
	TourGuideDiscountPercentage float32   `db:"tour_guide_discount_percentage"`
	Price                       int64     `db:"price"`
	DiscountPercentage          float32   `db:"discount_percentage"`
	CreatedAt                   time.Time `db:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at"`
	Bookings                    []TourGuideBookings
}

type Locals struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Address     string    `db:"address"`
	City        string    `db:"city"`
	Province    string    `db:"province"`
	Longitude   string    `db:"longitude"`
	Latitude    string    `db:"latitude"`
	Label       string    `db:"label"`
	OpenedTime  string    `db:"opened_time"`
	PhotoUrl    string    `db:"photo_url"`
	IsBusiness  bool      `db:"is_business"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Reviews     []Review
}

type Review struct {
	ID        uuid.UUID `db:"id"`
	Star      int       `db:"star"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	PhotoURL  string    `db:"photo_url"`
	UserID    uuid.UUID `db:"user_id"`
	LocalID   uuid.UUID `db:"local_id"`
}
