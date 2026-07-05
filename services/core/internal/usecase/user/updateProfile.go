package user

import (
	"context"
	"core_service/internal/domain"
	"strings"
)

func (s *userService) UpdateProfile(ctx context.Context, upd UpdateProfileInput) error {
	if upd.Email == nil && upd.FirstName == nil && upd.LastName == nil {
		return ErrNoContent
	}

	if upd.Email != nil {
		if strings.TrimSpace(*upd.Email) == "" || !isValidEmail(*upd.Email) {
			return ErrInvalidEmail
		}
	}

	if upd.FirstName != nil && strings.TrimSpace(*upd.FirstName) == "" {
		return ErrInvalidName
	}

	if upd.LastName != nil && strings.TrimSpace(*upd.LastName) == "" {
		return ErrInvalidName
	}

	if err := s.userRepo.UpdateProfile(ctx, upd.UserID, &domain.UpdateUserProfile{
		Email:     upd.Email,
		FirstName: upd.FirstName,
		LastName:  upd.LastName,
	}); err != nil {
		return err
	}

	return nil
}
