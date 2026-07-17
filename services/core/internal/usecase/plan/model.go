package plan

import "github.com/google/uuid"

type CreatePlanInput struct {
	EmployeeID  uuid.UUID
	CreatedBy   uuid.UUID
	Title       string
	Description *string
}

type UpdatePlanInput struct {
	PlanID      uuid.UUID
	ManagerID   uuid.UUID
	Title       string
	Description *string
}

type CreateAIInput struct {
	EmployeeID  uuid.UUID
	Topic       string
	Description string
	TargetRole  string
	CreatedBy   uuid.UUID
}
