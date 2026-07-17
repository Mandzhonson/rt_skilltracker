package model

import (
	"time"

	"github.com/google/uuid"
)

type SkillModel struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Category    string    `db:"category"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type UserSkillModel struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	SkillID     uuid.UUID `db:"skill_id"`
	PlanID      uuid.UUID `db:"plan_id"`
	ConfirmedAt time.Time `db:"confirmed_at"`
}
