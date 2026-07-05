package auth

import (
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/postgres"
	"core_service/internal/repository/redis"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type authService struct {
	authRepo postgres.AuthRepository
	userRepo postgres.UserRepository
	jwt      *jwt.JWTService
	redis    redis.SessionRepository
}

func NewAuthService(authRepo postgres.AuthRepository, userRepo postgres.UserRepository, jwt *jwt.JWTService, r redis.SessionRepository) *authService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
		jwt:      jwt,
		redis:    r,
	}
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
