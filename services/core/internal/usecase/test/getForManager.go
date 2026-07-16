package test

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/usecase/user"

	"github.com/google/uuid"
)

func (s *testService) GetForManager(ctx context.Context, managerID uuid.UUID, planID uuid.UUID) (*domain.TestView, error) {

	if managerID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if planID == uuid.Nil {
		return nil, ErrInvalidPlanID
	}

	allowed, err := s.planRepo.ManagerOwnsPlan(ctx, managerID, planID)
	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, user.ErrForbidden
	}

	return s.getTestView(ctx, planID)
}
