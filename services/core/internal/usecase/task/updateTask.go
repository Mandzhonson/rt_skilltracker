package task

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func (s *taskService) Update(ctx context.Context, input UpdateTaskInput) error {

	if input.TaskID == uuid.Nil {
		return ErrInvalidTaskID
	}

	if input.Title == nil && input.Description == nil {
		return ErrInvalidUpdate
	}

	taskEntity, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	planEntity, err := s.planRepo.GetByID(ctx, taskEntity.PlanID)
	if err != nil {
		return err
	}

	if planEntity.CreatedBy != input.ManagerID {
		return ErrManagerForbidden
	}

	if input.Title != nil {
		if strings.TrimSpace(*input.Title) == "" {
			return ErrInvalidTitle
		}

	}

	err = s.taskRepo.Update(ctx, input.TaskID, input.Title, input.Description)

	if err != nil {

		if errors.Is(
			err,
			postgres.ErrTaskNotFound,
		) {
			return ErrTaskNotFound
		}

		return err
	}

	return nil
}
