package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"core_service/internal/repository/redis"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	hash := hashRefreshToken(refreshToken)

	session, err := s.redis.GetSession(ctx, claims.ID)
	if err != nil {
		if !errors.Is(err, redis.ErrSessionNotFound) {
			return "", "", err
		}
		token, err := s.authRepo.GetRefreshToken(ctx, claims.ID)
		if err != nil {
			if errors.Is(err, postgres.ErrRefreshTokenNotFound) {
				return "", "", ErrInvalidRefreshToken
			}
			return "", "", err
		}
		if token.Revoked {
			return "", "", ErrInvalidRefreshToken
		}
		if time.Now().After(token.ExpiresAt) {
			return "", "", ErrInvalidRefreshToken
		}
		if token.TokenHash != hash {
			return "", "", ErrInvalidRefreshToken
		}
		if err := s.redis.SaveSession(
			ctx,
			token.JTI,
			token.UserID,
			token.TokenHash,
			time.Until(token.ExpiresAt),
		); err != nil {
			return "", "", err
		}
		session = &redis.Session{
			UserID:    token.UserID,
			TokenHash: token.TokenHash,
		}
	}

	if session.TokenHash != hash {
		return "", "", ErrInvalidRefreshToken
	}
	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return "", "", err
	}
	newAccess, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}
	newRefresh, newJTI, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}
	newHash := hashRefreshToken(newRefresh)
	if err := s.authRepo.DeleteRefreshToken(ctx, claims.ID); err != nil {
		return "", "", err
	}
	if err := s.redis.DeleteSession(ctx, claims.ID); err != nil {
		return "", "", err
	}
	entity := domain.RefreshToken{
		UserID:    user.ID,
		JTI:       newJTI,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(s.jwt.RefreshTTL()),
	}
	if err := s.authRepo.SaveRefreshToken(ctx, entity); err != nil {
		return "", "", err
	}
	if err := s.redis.SaveSession(
		ctx,
		newJTI,
		user.ID,
		newHash,
		s.jwt.RefreshTTL(),
	); err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}
