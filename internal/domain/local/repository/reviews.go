package repository

import (
	"context"

	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

// GetReviewsByLocalBusinessID retrieves all reviews for a specific local business
func (r *localRepository) GetReviewsByLocalBusinessID(ctx context.Context, localBusinessID string, out *[]local.Review) error {
	query := `
		SELECT 
			r.id, r.star, r.content, r.created_at, r.updated_at, 
			u.photo_url, u.full_name as user_name
		FROM reviews r 
		INNER JOIN users u ON u.id = r.user_id
		WHERE r.local_id = $1
		ORDER BY r.created_at DESC`

	rows, err := r.queryExecutor.QueryxContext(ctx, query, localBusinessID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var result []local.Review
	for rows.Next() {
		var review local.Review
		if err := rows.StructScan(&review); err != nil {
			return err
		}
		result = append(result, review)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*out = result
	return nil
}
