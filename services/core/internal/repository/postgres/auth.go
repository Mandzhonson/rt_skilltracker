package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type authRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *authRepository {
	return &authRepository{
		pool: pool,
	}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, e domain.RefreshToken) error {
	m := converter.ToRefreshTokenModel(e)
	query := `INSERT INTO refresh_tokens(user_id, jti, token_hash, expires_at, revoked, created_at) VALUES($1,$2,$3,$4,$5,$6)`
	if _, err := r.pool.Exec(ctx, query, m.UserID, m.JTI, m.TokenHash, m.ExpiresAt, m.Revoked, m.CreatedAt); err != nil {
		return fmt.Errorf("repository.SaveRefreshToken: %w", err)
	}
	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshToken, error) {
	var m model.RefreshTokenModel

	query := `
	SELECT user_id, jti, token_hash, expires_at, revoked, created_at
	FROM refresh_tokens
	WHERE jti = $1  AND revoked = FALSE`

	err := r.pool.QueryRow(ctx, query, jti).Scan(
		&m.UserID,
		&m.JTI,
		&m.TokenHash,
		&m.ExpiresAt,
		&m.Revoked,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("repository.GetRefreshToken: %w", err)
	}

	return converter.ToRefreshTokenEntity(&m), nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, jti string) error {
	query := `
	UPDATE refresh_tokens
	SET revoked = TRUE
	WHERE jti = $1`

	_, err := r.pool.Exec(ctx, query, jti)
	if err != nil {
		return fmt.Errorf("repository.DeleteRefreshToken: %w", err)
	}

	return nil
}
