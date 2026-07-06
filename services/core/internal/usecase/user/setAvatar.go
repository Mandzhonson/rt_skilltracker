package user

import (
	"context"
)

const maxAvatarSize = 5 * 1024 * 1024

func (s *userService) SetAvatar(ctx context.Context, input SetAvatarInput) error {
	if input.Size > maxAvatarSize {
		return ErrAvatarTooLarge
	}

	extension, err := avatarExtension(input.ContentType)
	if err != nil {
		return err
	}

	userEntity, err := s.userRepo.GetById(ctx, input.UserID)
	if err != nil {
		return err
	}

	objectName := input.UserID.String() + extension

	objectKey, err := s.storage.UploadAvatar(ctx, objectName, input.File, input.Size, input.ContentType)
	if err != nil {
		return err
	}

	oldAvatar := userEntity.AvatarKey

	if err := s.userRepo.UpdateAvatar(ctx, input.UserID, &objectKey); err != nil {
		_ = s.storage.DeleteAvatar(ctx, objectKey)
		return err
	}

	if oldAvatar != nil && *oldAvatar != objectKey {
		_ = s.storage.DeleteAvatar(ctx, *oldAvatar)
	}

	return nil
}

func avatarExtension(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/webp":
		return ".webp", nil
	default:
		return "", ErrInvalidAvatarFormat
	}
}
