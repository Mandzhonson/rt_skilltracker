package task

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (s *taskService) Create(ctx context.Context, input CreateTaskInput) (uuid.UUID, error) {

	if strings.TrimSpace(input.Title) == "" {
		return uuid.Nil, ErrInvalidTitle
	}

	planEntity, err := s.planRepo.GetByID(
		ctx,
		input.PlanID,
	)

	if err != nil {

		if errors.Is(err, postgres.ErrPlanNotFound) {
			return uuid.Nil, ErrPlanNotFound
		}

		return uuid.Nil, err
	}

	if planEntity.CreatedBy != input.CreatedBy {
		return uuid.Nil, ErrForbidden
	}

	position, err := s.taskRepo.GetNextPosition(
		ctx,
		input.PlanID,
	)

	if err != nil {
		return uuid.Nil, err
	}

	entity := domain.NewTask(
		input.PlanID,
		input.Title,
		input.Description,
		position,
	)

	id, err := s.taskRepo.Create(ctx, entity)

	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"create task: %w",
			err,
		)
	}

	return id, nil
}
