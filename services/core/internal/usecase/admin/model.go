package admin

import (
	"core_service/internal/domain"

	"github.com/google/uuid"
)

type ListUsersInput struct {
	Page   int
	Limit  int
	Role   *domain.Role
	Search *string
}

type UpdateRoleInput struct {
	ActorID uuid.UUID
	UserID  uuid.UUID
	Role    domain.Role
}

type AssignManagerInput struct {
	UserID    uuid.UUID
	ManagerID uuid.UUID
}
