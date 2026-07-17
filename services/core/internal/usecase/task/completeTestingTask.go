package task

import (
	"context"

	"github.com/google/uuid"
)

func (s *TaskService) CompleteTestingTask(
	ctx context.Context,
	planID uuid.UUID,
	userID uuid.UUID,
) error {

	plan, err := s.planRepo.GetByID(
		ctx,
		planID,
	)

	if err != nil {
		return err
	}

	if plan.EmployeeID != userID {
		return ErrEmployeeForbidden
	}

	err = s.taskRepo.CompleteTestingTask(
		ctx,
		planID,
	)

	if err != nil {
		return err
	}

	_, err = s.planRepo.RecalculateProgress(
		ctx,
		planID,
	)

	return err
}
