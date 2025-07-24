package session

import (
	"time"

	"github.com/google/uuid"
)

type Table struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
