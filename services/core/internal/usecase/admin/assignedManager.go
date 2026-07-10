package admin

import (
	"context"
	"core_service/internal/domain"
	"errors"
)

var (
	ErrInvalidManager = errors.New("user is not manager")
	ErrAssignYourself = errors.New("cannot assign yourself")
	ErrManagerCycle   = errors.New("manager hierarchy contains cycle")
)

func (s *adminService) AssignManager(ctx context.Context, input AssignManagerInput) error {

	if input.UserID == input.ManagerID {
		return ErrAssignYourself
	}

	user, err := s.userRepo.GetById(ctx, input.UserID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleAdmin {
		return ErrInvalidManager
	}

	manager, err := s.userRepo.GetById(ctx, input.ManagerID)
	if err != nil {
		return err
	}

	if manager.Role != domain.RoleManager {
		return ErrInvalidManager
	}
	err = s.validateManagerHierarchy(ctx, input.UserID, input.ManagerID)
	if err != nil {
		return err
	}
	return s.userRepo.AssignManager(ctx, input.UserID, input.ManagerID)
}
