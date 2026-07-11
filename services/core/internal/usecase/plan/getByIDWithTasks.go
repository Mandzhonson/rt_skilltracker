package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (s *planService) GetByIDWithTasks(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.PlanWithTasks, error) {
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrPlanNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("get plan by id: %w", err)
	}

	if plan.CreatedBy != managerID {
		return nil, ErrManagerForbidden
	}

	planWithTasks, err := s.planRepo.GetByIDWithTasks(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrPlanNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("get plan with tasks: %w", err)
	}

	return planWithTasks, nil
}
