package task

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

func (s *TaskService) GetByID(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.Task, error) {

	if id == uuid.Nil {
		return nil, ErrInvalidTaskID
	}

	entity, err := s.taskRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {

		if errors.Is(err, postgres.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	planEntity, err := s.planRepo.GetByID(
		ctx,
		entity.PlanID,
	)

	if err != nil {
		return nil, err
	}

	if planEntity.CreatedBy != managerID {
		return nil, ErrManagerForbidden
	}

	return entity, nil
}
