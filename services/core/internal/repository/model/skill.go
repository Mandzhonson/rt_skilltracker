package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSkillModel struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	PlanID      uuid.UUID `db:"plan_id"`
	Name        string    `db:"name"`
	ConfirmedAt time.Time `db:"confirmed_at"`
}

