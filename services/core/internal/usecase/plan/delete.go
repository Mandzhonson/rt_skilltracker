package plan

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (s *planService) Delete(ctx context.Context, managerID uuid.UUID, planID uuid.UUID) error {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		if errors.Is(err, postgres.ErrPlanNotFound) {
			return ErrPlanNotFound
		}
		return fmt.Errorf("get plan: %w", err)
	}

	if plan.CreatedBy != managerID {
		return ErrManagerForbidden
	}

	err = s.planRepo.Delete(ctx, planID)
	if err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}

	return nil
}
