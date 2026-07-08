package task

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

func (s *taskService) Delete(ctx context.Context, managerID uuid.UUID, taskID uuid.UUID) error {
	if taskID == uuid.Nil {
		return ErrInvalidTaskID
	}

	entity, err := s.taskRepo.GetByID(ctx, taskID)

	if err != nil {

		if errors.Is(err, postgres.ErrTaskNotFound) {
			return ErrTaskNotFound
		}

		return err
	}

	planEntity, err := s.planRepo.GetByID(ctx, entity.PlanID)
	if err != nil {
		return err
	}

	if planEntity.CreatedBy != managerID {
		return ErrForbidden
	}

	err = s.taskRepo.Delete(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}
