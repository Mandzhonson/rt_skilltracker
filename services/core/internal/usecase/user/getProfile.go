package user

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
