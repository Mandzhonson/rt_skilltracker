package model

import (
	"time"

	"github.com/google/uuid"
)

type PlanModel struct {
	ID               uuid.UUID `db:"id"`
	EmployeeID       uuid.UUID `db:"employee_id"`
	CreatedBy        uuid.UUID `db:"created_by"`
	Title            string    `db:"title"`
	Description      *string   `db:"description"`
	CreationType     string    `db:"creation_type"`
	Progress         int16     `db:"progress"`
	GenerationStatus string    `db:"generation_status"`
	Status           string    `db:"status"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
