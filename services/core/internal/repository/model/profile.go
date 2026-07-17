package model

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeProfileModel struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	MiddleName *string
	Email      string
	Position   string
	AvatarPath *string
	CreatedAt  time.Time
}
