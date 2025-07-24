package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetAllTouristAttractions retrieves all tourist attractions with optional city filtering
func (r *localRepository) GetAllTouristAttractions(ctx context.Context, city string, out *[]local.TouristAttractions) error {
	query := `
		SELECT 
			id, name, description, address, city, province, longitude, latitude, 
			photo_url, tour_guide_price, tour_guide_count, tour_guide_discount_percentage, 
			price, discount_percentage, created_at, updated_at
		FROM tourist_attractions 
		WHERE 1=1`

	var rows *sqlx.Rows
	var err error

	if city != "" {
		query += " AND LOWER(city) LIKE $1"
		cityFilter := "%" + strings.ToLower(city) + "%"
		rows, err = r.queryExecutor.QueryxContext(ctx, query, cityFilter)
	} else {
		rows, err = r.queryExecutor.QueryxContext(ctx, query)
	}

	if err != nil {
		return err
	}
	defer rows.Close()

	var result []local.TouristAttractions
	for rows.Next() {
		var attraction local.TouristAttractions
		if err := rows.StructScan(&attraction); err != nil {
			return err
		}
		result = append(result, attraction)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*out = result
	return nil
}

// GetTouristAttractionByID retrieves a tourist attraction by its ID
func (r *localRepository) GetTouristAttractionByID(ctx context.Context, data *local.TouristAttractions) error {
	query := `
		SELECT 
			id, name, description, address, city, province, longitude, latitude, 
			photo_url, tour_guide_price, tour_guide_count, tour_guide_discount_percentage, 
			price, discount_percentage, created_at, updated_at
		FROM tourist_attractions
		WHERE id = $1`

	row := r.queryExecutor.QueryRowxContext(ctx, query, data.ID)
	if err := row.StructScan(data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return local.ErrLBNotFound
		}
		return err
	}

	return nil
}

// GetBookingsByTouristAttractionID retrieves all bookings for a specific tourist attraction
func (r *localRepository) GetBookingsByTouristAttractionID(ctx context.Context, attractionID string, out *[]local.TourGuideBookings) error {
	query := `
		SELECT 
			tb.id, tb.payment_url, tb.star, tb.content, tb.created_at, tb.updated_at, 
			tb.status, tb.user_id, tb.tourist_attraction_id
		FROM tourguide_bookings tb 
		INNER JOIN users u ON u.id = tb.user_id
		WHERE tb.tourist_attraction_id = $1
		ORDER BY tb.created_at DESC`

	rows, err := r.queryExecutor.QueryxContext(ctx, query, attractionID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var result []local.TourGuideBookings
	for rows.Next() {
		var booking local.TourGuideBookings
		if err := rows.StructScan(&booking); err != nil {
			return err
		}
		result = append(result, booking)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*out = result
	return nil
}

// GetFullyBookedDates retrieves dates when tourist attraction is fully booked for a specific month and year
func (r *localRepository) GetFullyBookedDates(ctx context.Context, attractionID string, year, month int, dates *[]string) error {
	query := `
		SELECT 
			DATE(tb.booked_at) AS date
		FROM tourguide_bookings tb
		JOIN tourist_attractions ta ON ta.id = tb.tourist_attraction_id
		WHERE EXTRACT(MONTH FROM tb.booked_at) = $2
			AND EXTRACT(YEAR FROM tb.booked_at) = $3
			AND tb.tourist_attraction_id = $1
			AND tb.status = 'confirmed'
		GROUP BY DATE(tb.booked_at), ta.tour_guide_count
		HAVING COUNT(*) >= ta.tour_guide_count
		ORDER BY DATE(tb.booked_at)`

	rows, err := r.queryExecutor.QueryxContext(ctx, query, attractionID, month, year)
	if err != nil {
		return err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return err
		}
		results = append(results, date)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*dates = results
	return nil
}

// CreateTourGuideBooking creates a new tour guide booking
func (r *localRepository) CreateTourGuideBooking(ctx context.Context, booking *local.TourGuideBookings) error {
	query := `
		INSERT INTO tourguide_bookings (
			id, payment_url, star, content, booked_at, status, user_id, tourist_attraction_id,
			created_at, updated_at
		) VALUES (
			:id, :payment_url, :star, :content, :booked_at, :status, :user_id, :tourist_attraction_id,
			NOW(), NOW()
		)`

	_, err := r.queryExecutor.NamedExecContext(ctx, query, booking)
	if err != nil {
		return err
	}

	return nil
}

// CreateTouristAttraction creates a new tourist attraction in the database
func (r *localRepository) CreateTouristAttraction(ctx context.Context, attraction *local.TouristAttractions) error {
	query := `
		INSERT INTO tourist_attractions (
			id, name, description, address, city, province, longitude, latitude,
			photo_url, tour_guide_price, tour_guide_count, tour_guide_discount_percentage,
			price, discount_percentage, created_at, updated_at
		) VALUES (
			:id, :name, :description, :address, :city, :province, :longitude, :latitude,
			:photo_url, :tour_guide_price, :tour_guide_count, :tour_guide_discount_percentage,
			:price, :discount_percentage, :created_at, :updated_at
		)`

	_, err := r.queryExecutor.NamedExecContext(ctx, query, attraction)
	return err
}

// UpdateTouristAttraction updates an existing tourist attraction in the database
func (r *localRepository) UpdateTouristAttraction(ctx context.Context, attraction *local.TouristAttractions) error {
	query := `
		UPDATE tourist_attractions SET
			name = :name,
			description = :description,
			address = :address,
			city = :city,
			province = :province,
			longitude = :longitude,
			latitude = :latitude,
			photo_url = :photo_url,
			tour_guide_price = :tour_guide_price,
			tour_guide_count = :tour_guide_count,
			tour_guide_discount_percentage = :tour_guide_discount_percentage,
			price = :price,
			discount_percentage = :discount_percentage,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.queryExecutor.NamedExecContext(ctx, query, attraction)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return local.ErrLBNotFound
	}

	return nil
}

// DeleteTouristAttraction deletes a tourist attraction from the database
func (r *localRepository) DeleteTouristAttraction(ctx context.Context, attractionID string) error {
	query := `DELETE FROM tourist_attractions WHERE id = $1`

	result, err := r.queryExecutor.ExecContext(ctx, query, attractionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return local.ErrLBNotFound
	}

	return nil
}
