package user

import (
	"core_service/internal/repository/postgres"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPassword   = errors.New("password must be at least 8 characters long")
	ErrInvalidName       = errors.New("first name and last name are required")
	ErrUserNotFound      = errors.New("user not found")
)

type userService struct {
	userRepo postgres.UserRepository
}

func NewUserService(userRepository postgres.UserRepository) *userService {
	return &userService{
		userRepo: userRepository,
	}
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
