package domain

import (
	"time"

	"github.com/google/uuid"
)

type PlanStatus string

const (
	PlanActive    PlanStatus = "active"
	PlanCompleted PlanStatus = "completed"
	PlanArchived  PlanStatus = "archived"
)

type CreationType string

const (
	CreationManual CreationType = "manual"
	CreationAI     CreationType = "ai"
)

type Plan struct {
	ID           uuid.UUID
	EmployeeID   uuid.UUID
	CreatedBy    uuid.UUID
	Title        string
	Description  *string
	CreationType CreationType
	Progress     int
	Status       PlanStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewPlan(employeeID uuid.UUID, createdBy uuid.UUID, title string, description *string, creationType CreationType) *Plan {
	return &Plan{
		EmployeeID:   employeeID,
		CreatedBy:    createdBy,
		Title:        title,
		Description:  description,
		CreationType: creationType,
		Progress:     0,
		Status:       PlanActive,
	}
}

type PlanWithTasks struct {
	Plan  *Plan
	Tasks []*Task
}
