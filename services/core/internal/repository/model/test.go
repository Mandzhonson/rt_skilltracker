package model

import (
	"time"

	"github.com/google/uuid"
)

type TestModel struct {
	ID        uuid.UUID
	PlanID    uuid.UUID
	CreatedAt time.Time
}

