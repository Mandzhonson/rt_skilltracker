package user

import (
	"core_service/internal/repository/minio"
	"core_service/internal/repository/postgres"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists   = errors.New("user with this email already exists")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidPassword     = errors.New("password must be at least 8 characters long")
	ErrInvalidName         = errors.New("first name and last name are required")
	ErrNoContent           = errors.New("no content")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidAvatarFormat = errors.New("invalid avatar format")
	ErrAvatarTooLarge      = errors.New("avatar is too large")
	ErrAvatarNotFound      = errors.New("avatar is not found")
	ErrNotManager          = errors.New("user is not a manager")
)

type userService struct {
	userRepo postgres.UserRepository
	storage  minio.Storage
}

func NewUserService(userRepository postgres.UserRepository, storage minio.Storage) *userService {
	return &userService{
		userRepo: userRepository,
		storage:  storage,
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
