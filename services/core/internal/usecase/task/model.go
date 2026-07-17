package task

import (
	"core_service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type CreateTaskInput struct {
	PlanID      uuid.UUID
	CreatedBy   uuid.UUID
	Title       string
	Description *string
}

type UpdateTaskInput struct {
	TaskID      uuid.UUID
	ManagerID   uuid.UUID
	Title       *string
	Description *string
}

type UpdateTaskStatusInput struct {
	TaskID uuid.UUID
	UserID uuid.UUID
	Status domain.TaskStatus
}

type EmployeePlanResponse struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description *string                `json:"description"`
	Progress    int                    `json:"progress"`
	Status      string                 `json:"status"`
	Tasks       []EmployeeTaskResponse `json:"tasks"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type EmployeeTaskResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Position    int     `json:"position"`
	Status      string  `json:"status"`
}
