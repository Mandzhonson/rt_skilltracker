package test

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *testService) GetForEmployee(ctx context.Context, userID uuid.UUID, planID uuid.UUID) (*domain.TestView, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return s.getTestView(ctx, planID)
}
