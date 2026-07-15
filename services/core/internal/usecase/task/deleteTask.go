package task

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"

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
		return ErrManagerForbidden
	}

	err = s.taskRepo.Delete(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	_, err = s.planRepo.RecalculateProgress(ctx, entity.PlanID)
	if err != nil {
		return fmt.Errorf("recalculate progress: %w", err)
	}

	return nil
}
