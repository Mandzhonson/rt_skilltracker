package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
	ManagerID    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
