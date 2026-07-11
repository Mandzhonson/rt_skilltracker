package user

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (s *userService) GetEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.User, error) {
	user, err := s.userRepo.GetById(ctx, managerID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user.Role != domain.RoleManager {
		return nil, ErrNotManager
	}

	employees, err := s.userRepo.ListEmployeesByManager(ctx, managerID)
	if err != nil {
		return nil, fmt.Errorf("list employees by manager: %w", err)
	}

	return employees, nil
}
