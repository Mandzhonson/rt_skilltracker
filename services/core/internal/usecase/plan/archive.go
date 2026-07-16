package plan

import (
	"context"

	"github.com/google/uuid"
)

func (s *planService) Archive(ctx context.Context, managerID, planID uuid.UUID) error {

	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	if plan.CreatedBy != managerID {
		return ErrForbidden
	}

	return s.planRepo.Archive(ctx, planID)
}
