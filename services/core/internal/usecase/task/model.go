package task

import "github.com/google/uuid"

type CreateTaskInput struct {
	PlanID      uuid.UUID
	CreatedBy   uuid.UUID
	Title       string
	Description *string
}
