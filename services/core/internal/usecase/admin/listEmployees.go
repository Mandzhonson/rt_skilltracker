package admin

import (
	"context"

	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *adminService) ListEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.User, error) {
	manager, err := s.userRepo.GetById(ctx, managerID)
	if err != nil {
		return nil, err
	}

	if manager.Role != domain.RoleManager {
		return nil, ErrInvalidManager
	}

	return s.userRepo.ListEmployeesByManager(ctx, managerID)
}
