package user

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (uuid.UUID, error) {
	u := domain.NewEmployee(
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
	)
	role := input.Role

	if role == "" {
		role = domain.RoleEmployee
	}
	if err := validateUser(u); err != nil {
		return uuid.Nil, err
	}

	_, err := s.userRepo.GetByEmail(ctx, u.Email)
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
	u.Role = role

	id, err := s.userRepo.Create(ctx, u)
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
	if strings.TrimSpace(u.FirstName) == "" || strings.TrimSpace(u.LastName) == "" {
		return ErrInvalidName
	}
	return nil
}
