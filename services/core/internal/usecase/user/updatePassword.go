package user

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return ErrInvalidPassword
	}
	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashNewPass, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.userRepo.UpdatePassword(ctx, userID, hashNewPass); err != nil {
		return err
	}
	return nil
}
