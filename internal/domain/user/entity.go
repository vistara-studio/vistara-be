package user

import (
	"time"

	"github.com/google/uuid"
)

type Table struct {
	ID           uuid.UUID    `db:"id"`
	FullName     string       `db:"full_name"`
	Email        string       `db:"email"`
	Password     string       `db:"password"`
	AuthProvider AuthProvider `db:"auth_provider"`
	PhotoUrl     string       `db:"photo_url"`
	IsPremium    bool         `db:"is_premium"`
	ExpiredAt    time.Time    `db:"expired_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}
