package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshTokenModel struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	JTI       string    `db:"jti"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}
