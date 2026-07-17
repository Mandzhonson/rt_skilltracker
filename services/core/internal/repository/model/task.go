package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskModel struct {
	ID          uuid.UUID `db:"id"`
	PlanID      uuid.UUID `db:"plan_id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
	Position    int16     `db:"position"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
