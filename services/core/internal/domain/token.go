package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    uuid.UUID
	JTI       string
	TokenHash string
	ExpiresAt time.Time
}
