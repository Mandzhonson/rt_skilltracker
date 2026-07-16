package test

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *testService) GetForAdmin(ctx context.Context, planID uuid.UUID) (*domain.TestView, error) {
	return s.getTestView(ctx, planID)
}
