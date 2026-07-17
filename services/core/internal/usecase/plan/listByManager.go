package plan

import (
	"context"
	"core_service/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (s *planService) ListByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.Plan, error) {
	plans, err := s.planRepo.ListByManager(ctx, managerID)
	if err != nil {
		return nil, fmt.Errorf("list plans by manager: %w", err)
	}
	return plans, nil
}
