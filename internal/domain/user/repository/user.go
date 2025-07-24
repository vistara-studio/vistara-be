package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vistara-studio/vistara-be/internal/domain/user"
	"github.com/lib/pq"
)

func (r *userRepository) CreateUser(ctx context.Context, data user.Table) error {
	query := `INSERT INTO users (
		id, full_name, email, password, auth_provider, photo_url
	) VALUES (
		:id, :full_name, :email, :password, :auth_provider, :photo_url
	)`

	_, err := r.q.NamedExecContext(ctx, query, data)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" && pqErr.Constraint == "accounts_email_key" {
				return user.ErrEmailAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *userRepository) GetAccountByEmail(ctx context.Context, data *user.Table) error {
	query := `SELECT 
	id, full_name, email, password, auth_provider, photo_url, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	row := r.q.QueryRowxContext(ctx, query, data.Email)
	if err := row.StructScan(data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ErrUserNotFound
		} else {
			return err
		}
	}

	return nil
}
