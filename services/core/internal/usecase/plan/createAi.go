package plan

import (
	"context"
	"errors"
	"strings"

	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"core_service/internal/usecase/user"

	"github.com/google/uuid"
)

func (s *planService) CreateAI(ctx context.Context, input CreateAIInput) (uuid.UUID, error) {
	if strings.TrimSpace(input.Topic) == "" {
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

	manager, err := s.userRepo.GetById(ctx, input.CreatedBy)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return uuid.Nil, user.ErrUserNotFound
		}
		return uuid.Nil, err
	}

	if !manager.IsManager() {
		return uuid.Nil, ErrInvalidCreator
	}

	if employee.ManagerID == nil || *employee.ManagerID != manager.ID {
		return uuid.Nil, ErrEmployeeNotAssigned
	}

	plan := domain.NewPlan(
		input.EmployeeID,
		input.CreatedBy,
		input.Topic,
		nil,
		domain.CreationAI,
	)

	plan.GenerationStatus = domain.GenerationPending

	planID, err := s.planRepo.Create(ctx, plan)
	if err != nil {
		return uuid.Nil, err
	}

	go s.generateAIPlan(context.Background(), planID)

	return planID, nil

}
