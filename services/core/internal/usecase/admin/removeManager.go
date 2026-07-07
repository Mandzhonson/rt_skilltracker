package admin

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrManagerNotAssigned = errors.New("manager is not assigned")

func (s *adminService) RemoveManager(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return err
	}

	if user.ManagerID == nil {
		return ErrManagerNotAssigned
	}

	return s.userRepo.RemoveManager(ctx, userID)
}
