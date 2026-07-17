package user

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"strings"
)

func (s *userService) UpdateProfile(ctx context.Context, upd UpdateProfileInput) (*domain.User, error) {

	if upd.Email == nil && upd.FirstName == nil && upd.LastName == nil {
		return nil, ErrNoContent
	}

	if upd.Email != nil {
		if strings.TrimSpace(*upd.Email) == "" || !isValidEmail(*upd.Email) {
			return nil, ErrInvalidEmail
		}
	}

	if upd.FirstName != nil && strings.TrimSpace(*upd.FirstName) == "" {
		return nil, ErrInvalidName
	}

	if upd.LastName != nil && strings.TrimSpace(*upd.LastName) == "" {
		return nil, ErrInvalidName
	}

	if err := s.userRepo.UpdateProfile(ctx, upd.UserID, &domain.UpdateUserProfile{
		Email:     upd.Email,
		FirstName: upd.FirstName,
		LastName:  upd.LastName,
	}); err != nil {

		switch {
		case errors.Is(err, postgres.ErrUserAlreadyExists):
			return nil, ErrUserAlreadyExists

		case errors.Is(err, postgres.ErrUserNotFound):
			return nil, ErrUserNotFound

		default:
			return nil, err
		}
	}

	return s.userRepo.GetById(ctx, upd.UserID)
}
