package plan

import "github.com/google/uuid"

type CreatePlanInput struct {
	EmployeeID  uuid.UUID
	CreatedBy   uuid.UUID
	Title       string
	Description *string
}
