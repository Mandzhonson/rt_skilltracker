package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {

	user, err := s.authenticateUser(ctx, email, password)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, jti, err := s.generateTokens(user)
	if err != nil {
		return "", "", err
	}

	if err := s.storeRefreshSession(ctx, user.ID, refreshToken, jti); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) authenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authService) generateTokens(user *domain.User) (accessToken string, refreshToken string, jti string, err error) {

	accessToken, err = s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, jti, err = s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, jti, nil
}

func (s *authService) storeRefreshSession(ctx context.Context, userID uuid.UUID, refreshToken string, jti string) error {

	hashToken := hashRefreshToken(refreshToken)

	now := time.Now().UTC()

	entity := domain.RefreshToken{
		UserID:    userID,
		JTI:       jti,
		TokenHash: hashToken,
		ExpiresAt: now.Add(s.jwt.RefreshTTL()),
	}

	if err := s.repo.SaveRefreshToken(ctx, entity); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}

	return nil
}
