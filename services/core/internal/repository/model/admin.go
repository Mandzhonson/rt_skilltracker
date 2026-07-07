package model

import "core_service/internal/domain"

type ListUsersParams struct {
	Offset int
	Limit  int
	Role   *domain.Role
	Search *string
}
