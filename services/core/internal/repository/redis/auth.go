package redis

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID    uuid.UUID
	TokenHash string
}

type SessionRepository interface {
	BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)

	SaveSession(ctx context.Context, refreshJTI string, userID uuid.UUID, tokenHash string, ttl time.Duration) error
	GetSession(ctx context.Context, refreshJTI string) (*Session, error)
	DeleteSession(ctx context.Context, refreshJTI string) error
}

type RedisSessionRepository struct {
	r *redis.Client
}

func NewRedisSessionRepository(r *redis.Client) *RedisSessionRepository {
	return &RedisSessionRepository{
		r: r,
	}
}

func (s *RedisSessionRepository) BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error {
	return s.r.Set(ctx, blacklistKey(jti), "1", ttl).Err()
}

func (s *RedisSessionRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	exists, err := s.r.Exists(ctx, blacklistKey(jti)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (s *RedisSessionRepository) SaveSession(
	ctx context.Context,
	refreshJTI string,
	userID uuid.UUID,
	tokenHash string,
	ttl time.Duration,
) error {

	key := sessionKey(refreshJTI)

	if err := s.r.HSet(
		ctx,
		key,
		"user_id", userID.String(),
		"token_hash", tokenHash,
	).Err(); err != nil {
		return err
	}

	return s.r.Expire(ctx, key, ttl).Err()
}

func (s *RedisSessionRepository) GetSession(
	ctx context.Context,
	refreshJTI string,
) (*Session, error) {

	values, err := s.r.HGetAll(ctx, sessionKey(refreshJTI)).Result()
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrSessionNotFound
	}

	userID, err := uuid.Parse(values["user_id"])
	if err != nil {
		return nil, err
	}

	return &Session{
		UserID:    userID,
		TokenHash: values["token_hash"],
	}, nil
}

func (s *RedisSessionRepository) DeleteSession(
	ctx context.Context,
	refreshJTI string,
) error {
	return s.r.Del(ctx, sessionKey(refreshJTI)).Err()
}

func blacklistKey(jti string) string {
	return "auth:blacklist:" + jti
}

func sessionKey(refreshJTI string) string {
	return "auth:session:" + refreshJTI
}
