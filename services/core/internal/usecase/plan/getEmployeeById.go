package plan

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *planService) GetEmployeePlan(ctx context.Context, employeeID, planID uuid.UUID) (*domain.PlanWithTasks, error) {

	if planID == uuid.Nil {
		return nil, ErrInvalidPlanID
	}

	planEntity, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if planEntity.EmployeeID != employeeID {
		return nil, ErrEmployeeForbidden
	}

	tasks, err := s.taskRepo.ListByPlanID(ctx, planID)
	if err != nil {
		return nil, err
	}

	return &domain.PlanWithTasks{
		Plan:  planEntity,
		Tasks: tasks,
	}, nil
}
