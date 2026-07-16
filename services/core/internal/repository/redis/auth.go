package redis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID    uuid.UUID
	TokenHash string
}

//go:generate mockgen -source=auth.go -destination=mocks/mock_redis.go -package=mocks
type SessionRepository interface {
	BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)

	SaveSession(ctx context.Context, refreshJTI string, userID uuid.UUID, tokenHash string, ttl time.Duration) error
	GetSession(ctx context.Context, refreshJTI string) (*Session, error)
	DeleteSession(ctx context.Context, refreshJTI string) error
}

type RedisSessionRepository struct {
	r   *redis.Client
	log *slog.Logger
}

func NewRedisSessionRepository(r *redis.Client, log *slog.Logger) *RedisSessionRepository {
	return &RedisSessionRepository{
		r:   r,
		log: log,
	}
}

func (s *RedisSessionRepository) BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error {
	err := s.r.Set(ctx, blacklistKey(jti), "1", ttl).Err()
	if err != nil {
		s.log.Error("Failed to blacklist access token",
			slog.String("error", err.Error()),
			slog.String("jti", jti),
			slog.Duration("ttl", ttl),
		)
		return err
	}
	return nil
}

func (s *RedisSessionRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	exists, err := s.r.Exists(ctx, blacklistKey(jti)).Result()
	if err != nil {
		s.log.Error("Failed to check if token is blacklisted",
			slog.String("error", err.Error()),
			slog.String("jti", jti),
		)
		return false, err
	}

	isBlacklisted := exists == 1
	if isBlacklisted {
		s.log.Info("Token is blacklisted",
			slog.String("jti", jti),
		)
	}
	return isBlacklisted, nil
}

func (s *RedisSessionRepository) SaveSession(ctx context.Context, refreshJTI string, userID uuid.UUID, tokenHash string, ttl time.Duration) error {

	key := sessionKey(refreshJTI)

	if err := s.r.HSet(ctx, key, "user_id", userID.String(), "token_hash", tokenHash).Err(); err != nil {
		s.log.Error("Failed to save session",
			slog.String("error", err.Error()),
			slog.String("refresh_jti", refreshJTI),
			slog.String("user_id", userID.String()),
			slog.Duration("ttl", ttl),
		)
		return err
	}

	if err := s.r.Expire(ctx, key, ttl).Err(); err != nil {
		s.log.Error("Failed to set session expiration",
			slog.String("error", err.Error()),
			slog.String("refresh_jti", refreshJTI),
			slog.String("user_id", userID.String()),
			slog.Duration("ttl", ttl),
		)
		return err
	}

	return nil
}

func (s *RedisSessionRepository) GetSession(
	ctx context.Context,
	refreshJTI string,
) (*Session, error) {

	values, err := s.r.HGetAll(ctx, sessionKey(refreshJTI)).Result()
	if err != nil {
		s.log.Error("Failed to get session",
			slog.String("error", err.Error()),
			slog.String("refresh_jti", refreshJTI),
		)
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrSessionNotFound
	}

	userID, err := uuid.Parse(values["user_id"])
	if err != nil {
		s.log.Error("Failed to parse user_id from session",
			slog.String("error", err.Error()),
			slog.String("refresh_jti", refreshJTI),
			slog.String("user_id_value", values["user_id"]),
		)
		return nil, err
	}

	return &Session{
		UserID:    userID,
		TokenHash: values["token_hash"],
	}, nil
}

func (s *RedisSessionRepository) DeleteSession(ctx context.Context, refreshJTI string) error {
	err := s.r.Del(ctx, sessionKey(refreshJTI)).Err()
	if err != nil {
		s.log.Error("Failed to delete session",
			slog.String("error", err.Error()),
			slog.String("refresh_jti", refreshJTI),
		)
		return err
	}

	return nil
}

func blacklistKey(jti string) string {
	return "auth:blacklist:" + jti
}

func sessionKey(refreshJTI string) string {
	return "auth:session:" + refreshJTI
}
