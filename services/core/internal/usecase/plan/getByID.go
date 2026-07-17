package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

func (s *planService) GetByID(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.Plan, error) {

	if id == uuid.Nil {
		return nil, ErrInvalidPlanID
	}

	entity, err := s.planRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {

		if errors.Is(err, postgres.ErrPlanNotFound) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	if entity.CreatedBy != managerID {
		return nil, ErrForbidden
	}

	return entity, nil
}
