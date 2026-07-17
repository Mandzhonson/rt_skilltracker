package plan

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *planService) ListEmployeePlans(ctx context.Context, employeeID uuid.UUID) ([]*domain.PlanWithTasks, error) {

	plans, err := s.planRepo.ListByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PlanWithTasks, 0, len(plans))
	for _, p := range plans {
		tasks, err := s.taskRepo.ListByPlanID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, &domain.PlanWithTasks{
			Plan:  p,
			Tasks: tasks,
		})
	}

	return result, nil
}
