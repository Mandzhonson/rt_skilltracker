package auth

import (
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/postgres"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters long")
	ErrInvalidName        = errors.New("first name and last name are required")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type authService struct {
	repo postgres.AuthRepository
	jwt  jwt.JWTService
}

func NewAuthService(repo postgres.AuthRepository, jwt jwt.JWTService) *authService {
	return &authService{
		repo: repo,
		jwt:  jwt,
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

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
