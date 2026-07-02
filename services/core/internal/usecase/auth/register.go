package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (s *authService) CreateUser(ctx context.Context, u *domain.User) (uuid.UUID, error) {
	if err := validateUser(u); err != nil {
		return uuid.Nil, err
	}

	_, err := s.repo.GetByEmail(ctx, u.Email)
	switch {
	case err == nil:
		return uuid.Nil, ErrUserAlreadyExists
	case errors.Is(err, postgres.ErrUserNotFound):
	default:
		return uuid.Nil, fmt.Errorf("check user exists: %w", err)
	}

	hashedPassword, err := hashPassword(u.PasswordHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	u.PasswordHash = hashedPassword

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func validateUser(u *domain.User) error {
	if u == nil {
		return errors.New("user cannot be nil")
	}
	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}
	if len(u.PasswordHash) < 8 {
		return ErrInvalidPassword
	}
	if u.FirstName == "" || u.LastName == "" {
		return ErrInvalidName
	}
	return nil
}
