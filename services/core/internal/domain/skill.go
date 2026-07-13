package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserSkill struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	PlanID      uuid.UUID
	Name        string
	ConfirmedAt time.Time
}
