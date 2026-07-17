package admin

import (
	"context"
	"core_service/internal/domain"
	"errors"
)

var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrChangeOwnRole      = errors.New("cannot change your own role")
	ErrLastAdminProtected = errors.New("cannot remove the last administrator")
)

func (s *adminService) UpdateRole(ctx context.Context, input UpdateRoleInput) error {
	
	user, err := s.userRepo.GetById(ctx, input.UserID)
	if err != nil {
		return err
	}
	if user.Role == domain.RoleManager &&
		input.Role != domain.RoleManager {
		err = s.userRepo.ClearManagerAssignments(ctx, user.ID)
		if err != nil {
			return err
		}

	}
	return s.userRepo.UpdateRole(ctx, input.UserID, input.Role)
}
