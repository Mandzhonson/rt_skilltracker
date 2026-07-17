package admin

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func (s *adminService) ListUsers(ctx context.Context, input ListUsersInput) ([]*domain.User, error) {
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	params := model.ListUsersParams{
		Offset: (input.Page - 1) * input.Limit,
		Limit:  input.Limit,
		Role:   input.Role,
		Search: input.Search,
	}

	return s.userRepo.ListUsers(ctx, params)
}
