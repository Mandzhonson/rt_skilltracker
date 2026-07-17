package user

import (
	"context"
	"io"

	"github.com/google/uuid"
)

func (s *userService) GetAvatar(ctx context.Context, userID uuid.UUID) (io.ReadCloser, string, error) {
	userEntity, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return nil, "", err
	}

	if userEntity.AvatarKey == nil {
		return nil, "", ErrAvatarNotFound
	}

	return s.storage.GetAvatar(ctx, *userEntity.AvatarKey)
}
