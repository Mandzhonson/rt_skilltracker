package user

import (
	"context"

	"github.com/google/uuid"
)

func (s *userService) DeleteAvatar(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return err
	}

	if user.AvatarKey == nil {
		return ErrNoContent
	}

	err = s.storage.DeleteAvatar(ctx, *user.AvatarKey)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateAvatar(ctx, userID, nil)
}
