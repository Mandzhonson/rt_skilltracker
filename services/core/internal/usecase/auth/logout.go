package auth

import (
	"context"
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/postgres"
	"core_service/internal/repository/redis"
	"errors"
	"time"
)

func (s *authService) Logout(ctx context.Context, claimsAccess *jwt.Claims, refreshToken string) error {

	claimsRefresh, err := s.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	if claimsAccess.Subject != claimsRefresh.Subject {
		return ErrInvalidCredentials
	}

	if err := s.repo.DeleteRefreshToken(ctx, claimsRefresh.ID); err != nil &&
		!errors.Is(err, postgres.ErrRefreshTokenNotFound) {
		return err
	}

	if err := s.redis.DeleteSession(ctx, claimsRefresh.ID); err != nil && !errors.Is(err, redis.ErrSessionNotFound) {
		return err
	}

	ttl := time.Until(claimsAccess.ExpiresAt.Time)
	if ttl > 0 {
		if err := s.redis.BlacklistAccessToken(ctx, claimsAccess.ID, ttl); err != nil {
			return err
		}
	}
	return nil
}
