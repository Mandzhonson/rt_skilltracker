package admin

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *adminService) AdminGetPlan(ctx context.Context, planID uuid.UUID) (*domain.PlanWithTasks, error) {

	planWithTasks, err := s.planRepo.GetByIDWithTasks(ctx, planID)
	if err != nil {
		return nil, err
	}

	return planWithTasks, nil
}
