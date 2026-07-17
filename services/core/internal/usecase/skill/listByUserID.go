package skill

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *skillService) ListByUserID(ctx context.Context, requesterID uuid.UUID, userID uuid.UUID) ([]*domain.Skill, error) {

	requester, err := s.userRepo.GetById(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	if requester.IsEmployee() && requester.ID != userID {
		return nil, ErrForbidden
	}

	if requester.IsManager() {

		employee, err := s.userRepo.GetById(ctx, userID)
		if err != nil {
			return nil, err
		}

		if employee.ManagerID == nil ||
			*employee.ManagerID != requester.ID {
			return nil, ErrForbidden
		}
	}

	return s.skillRepo.ListByUserID(ctx, userID)
}
