package plan

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"
	"strings"
)

func (s *planService) Update(ctx context.Context, input UpdatePlanInput) error {
	title := strings.TrimSpace(input.Title)

	if title == "" {
		return ErrInvalidTitle
	}

	if len(title) < 3 {
		return ErrInvalidTitle
	}

	plan, err := s.planRepo.GetByID(ctx, input.PlanID)
	if err != nil {
		if errors.Is(err, postgres.ErrPlanNotFound) {
			return ErrPlanNotFound
		}
		return fmt.Errorf("get plan: %w", err)
	}

	if plan.CreatedBy != input.ManagerID {
		return ErrManagerForbidden
	}

	err = s.planRepo.Update(ctx, input.PlanID, title, input.Description)
	if err != nil {
		return fmt.Errorf("update plan: %w", err)
	}

	return nil
}
