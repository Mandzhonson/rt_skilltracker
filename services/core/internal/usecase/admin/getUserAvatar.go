package admin

import (
	"context"
	"core_service/internal/usecase/user"
	"io"

	"github.com/google/uuid"
)

func (s *adminService) GetUserAvatar(ctx context.Context, userID uuid.UUID) (io.ReadCloser, string, error) {
	userEntity, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	if userEntity.AvatarKey == nil {
		return nil, "", user.ErrAvatarNotFound
	}

	return s.storage.GetAvatar(ctx, *userEntity.AvatarKey)
}
