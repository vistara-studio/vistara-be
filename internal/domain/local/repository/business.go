package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetAllLocalBusinesses retrieves all local businesses with optional city and type filtering
func (r *localRepository) GetAllLocalBusinesses(ctx context.Context, params local.QueryParamRequestGetLocals, out *[]local.Locals) error {
	query := `
		SELECT 
			id, name, description, address, city, province, longitude, latitude, 
			label, opened_time, photo_url, is_business, created_at, updated_at
		FROM locals
		WHERE 1=1`

	queryParams := make(map[string]interface{})

	if params.City != "" {
		query += " AND LOWER(city) LIKE :city"
		queryParams["city"] = "%" + strings.ToLower(params.City) + "%"
	}

	if params.Type == "business" {
		query += " AND is_business = true"
	} else if params.Type == "individual" {
		query += " AND is_business = false"
	}

	query += " ORDER BY created_at DESC"

	// If no city filter, use regular query
	if params.City == "" {
		rows, err := r.queryExecutor.QueryxContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		var result []local.Locals
		for rows.Next() {
			var business local.Locals
			if err := rows.StructScan(&business); err != nil {
				return err
			}
			result = append(result, business)
		}

		if err := rows.Err(); err != nil {
			return err
		}

		*out = result
		return nil
	}

	// Use named query for city filtering
	rows, err := r.queryExecutor.NamedQueryContext(ctx, query, queryParams)
	if err != nil {
		return err
	}
	defer rows.Close()

	var result []local.Locals
	for rows.Next() {
		var business local.Locals
		if err := rows.StructScan(&business); err != nil {
			return err
		}
		result = append(result, business)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*out = result
	return nil
}

// GetLocalBusinessByID retrieves a local business by its ID
func (r *localRepository) GetLocalBusinessByID(ctx context.Context, business *local.Locals) error {
	query := `
		SELECT 
			id, name, description, address, city, province, longitude, latitude, 
			label, opened_time, photo_url, is_business, created_at, updated_at
		FROM locals
		WHERE id = $1`

	row := r.queryExecutor.QueryRowxContext(ctx, query, business.ID)
	if err := row.StructScan(business); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return local.ErrLBNotFound
		}
		return err
	}

	return nil
}

// CreateLocalBusiness creates a new local business in the database
func (r *localRepository) CreateLocalBusiness(ctx context.Context, business *local.Locals) error {
	query := `
		INSERT INTO locals (
			id, name, description, address, city, province, longitude, latitude,
			label, opened_time, photo_url, is_business, created_at, updated_at
		) VALUES (
			:id, :name, :description, :address, :city, :province, :longitude, :latitude,
			:label, :opened_time, :photo_url, :is_business, :created_at, :updated_at
		)`

	_, err := r.queryExecutor.NamedExecContext(ctx, query, business)
	if err != nil {
		return err
	}

	return nil
}

// UpdateLocalBusiness updates an existing local business in the database
func (r *localRepository) UpdateLocalBusiness(ctx context.Context, business *local.Locals) error {
	query := `
		UPDATE locals SET
			name = :name,
			description = :description,
			address = :address,
			city = :city,
			province = :province,
			longitude = :longitude,
			latitude = :latitude,
			label = :label,
			opened_time = :opened_time,
			photo_url = :photo_url,
			is_business = :is_business,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.queryExecutor.NamedExecContext(ctx, query, business)
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

// DeleteLocalBusiness deletes a local business from the database
func (r *localRepository) DeleteLocalBusiness(ctx context.Context, businessID string) error {
	query := `DELETE FROM locals WHERE id = $1`

	result, err := r.queryExecutor.ExecContext(ctx, query, businessID)
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
