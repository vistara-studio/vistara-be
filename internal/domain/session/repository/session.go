package repository

import (
	"context"

	"github.com/vistara-studio/vistara-be/internal/domain/session"
	"github.com/vistara-studio/vistara-be/internal/domain/user"
)

func (r *sessionRepository) CreateSession(ctx context.Context, data session.Table) error {
	query := `INSERT INTO sessions (
		id, user_id
	) VALUES (
		:id, :user_id
	)`

	_, err := r.q.NamedExecContext(ctx, query, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) DeleteOldestSessionByUserID(ctx context.Context, data session.Table) error {
	query := `DELETE FROM sessions
WHERE id = (
    SELECT id FROM sessions
    WHERE user_id = $1
    ORDER BY created_at ASC
    LIMIT 1
)`

	res, err := r.q.ExecContext(ctx, query, data.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *sessionRepository) GetSessionByUserID(ctx context.Context, data *user.Table, out *[]session.Table) error {
	query := `SELECT 
	id, user_id, created_at
	FROM sessions
	WHERE user_id = $1
	`

	rows, err := r.q.QueryxContext(ctx, query, data.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var result []session.Table
	for rows.Next() {
		var item session.Table
		if err := rows.StructScan(&item); err != nil {
			return err
		}

		result = append(result, item)
	}

	*out = result
	return nil
}
