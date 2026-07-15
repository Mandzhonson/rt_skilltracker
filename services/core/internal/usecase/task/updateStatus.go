package task

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

func (s *taskService) UpdateStatus(ctx context.Context, input UpdateTaskStatusInput) error {

	if input.TaskID == uuid.Nil {
		return ErrInvalidTaskID
	}

	if input.UserID == uuid.Nil {
		return ErrInvalidUserID
	}

	switch input.Status {
	case domain.TaskTodo,
		domain.TaskInProgress,
		domain.TaskDone:
	default:
		return ErrInvalidStatus
	}

	taskEntity, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	if taskEntity.Status == input.Status {
		return ErrInvalidStatus
	}

	planEntity, err := s.planRepo.GetByID(ctx, taskEntity.PlanID)
	if err != nil {
		if errors.Is(err, postgres.ErrPlanNotFound) {
			return ErrPlanNotFound
		}
		return err
	}

	if planEntity.ID != taskEntity.PlanID {
		return ErrPlanNotFound
	}

	if planEntity.EmployeeID != input.UserID {
		return ErrEmployeeForbidden
	}

	err = s.taskRepo.UpdateStatus(ctx, input.TaskID, string(input.Status))
	if err != nil {
		return err
	}
	progress, err := s.planRepo.RecalculateProgress(ctx, taskEntity.PlanID)

	if err != nil {
		return err
	}

	if progress == 100 {
		go func(planID uuid.UUID) {
			err := s.planCompletionService.GenerateSkillsIfCompleted(context.Background(), planID)
			if err != nil {
				// TODO: добавить slog.Error()
			}
		}(taskEntity.PlanID)
	}
	return nil
}
