package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPassword   = errors.New("password must be at least 8 characters long")
	ErrInvalidName       = errors.New("first name and last name are required")
)


type authService struct {
	repo postgres.UserRepository
}

func NewAuthService(repo postgres.UserRepository) *authService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) CreateUser(ctx context.Context, u *domain.User) (uuid.UUID, error) {
	if err := validateUser(u); err != nil {
		return uuid.Nil, err
	}

	// existingUser, err := s.repo.GetByEmail(ctx, u.Email)
	// if err == nil && existingUser != nil {
	// 	return uuid.Nil, ErrUserAlreadyExists
	// }

	hashedPassword, err := hashPassword(u.PasswordHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = hashedPassword

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
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

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
