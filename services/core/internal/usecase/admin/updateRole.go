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
	switch input.Role {
	case domain.RoleAdmin,
		domain.RoleManager,
		domain.RoleEmployee:
	default:
		return ErrInvalidRole
	}
	if input.ActorID == input.UserID {
		return ErrChangeOwnRole
	}
	user, err := s.userRepo.GetById(ctx, input.UserID)
	if err != nil {
		return err
	}
	if user.Role == domain.RoleAdmin && input.Role != domain.RoleAdmin {
		count, err := s.userRepo.CountAdmins(ctx)
		if err != nil {
			return err
		}
		if count == 1 {
			return ErrLastAdminProtected
		}
	}
	return s.userRepo.UpdateRole(ctx, input.UserID, input.Role)
}
