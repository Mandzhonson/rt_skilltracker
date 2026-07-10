package admin

import (
	"context"

	"github.com/google/uuid"
)

func (s *adminService) validateManagerHierarchy(ctx context.Context, userID uuid.UUID, managerID uuid.UUID) error {
	current := managerID
	for {
		manager, err := s.userRepo.GetById(
			ctx,
			current,
		)
		if err != nil {
			return err
		}
		if manager.ManagerID == nil {
			return nil
		}
		if *manager.ManagerID == userID {
			return ErrManagerCycle
		}
		current = *manager.ManagerID
	}
}
