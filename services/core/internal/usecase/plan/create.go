package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"core_service/internal/usecase/user"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (s *planService) Create(ctx context.Context, input CreatePlanInput) (uuid.UUID, error) {

	if strings.TrimSpace(input.Title) == "" {
		return uuid.Nil, ErrInvalidTitle
	}

	employee, err := s.userRepo.GetById(ctx, input.EmployeeID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return uuid.Nil, user.ErrUserNotFound
		}

		return uuid.Nil, err
	}

	if !employee.IsEmployee() {
		return uuid.Nil, ErrInvalidEmployee
	}

	creator, err := s.userRepo.GetById(ctx, input.CreatedBy)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return uuid.Nil, user.ErrUserNotFound
		}

		return uuid.Nil, err
	}

	if !creator.IsManager() {
		return uuid.Nil, ErrInvalidCreator
	}

	if employee.ManagerID == nil ||
		*employee.ManagerID != creator.ID {

		return uuid.Nil, ErrEmployeeNotAssigned
	}

	entity := domain.NewPlan(
		input.EmployeeID,
		input.CreatedBy,
		input.Title,
		input.Description,
		domain.CreationManual,
	)

	entity.GenerationStatus = domain.GenerationPending

	id, err := s.planRepo.Create(ctx, entity)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create plan: %w", err)
	}

	entity.ID = id

	go s.generateManualPlan(context.Background(), entity)

	return id, nil
}
