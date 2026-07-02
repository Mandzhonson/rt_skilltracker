package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type authRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *authRepository {
	return &authRepository{
		pool: pool,
	}
}

func (r *authRepository) Create(ctx context.Context, u *domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	m := converter.ToUserModel(u)
	query := `INSERT INTO users(email, password_hash, first_name, last_name, manager_id) VALUES($1,$2,$3,$4,$5) RETURNING id`
	err := r.pool.QueryRow(ctx, query, m.Email, m.PasswordHash, m.FirstName, m.LastName, m.ManagerID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repository.Create (user): %w", err)
	}
	return id, nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m model.UserModel

	query := `
	SELECT id, email, password_hash, first_name, last_name, role, manager_id, created_at, updated_at
	FROM users
	WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.FirstName,
		&m.LastName,
		&m.Role,
		&m.ManagerID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("repository.GetByEmail(user): %w", err)
	}

	return converter.ToUserEntity(&m), nil
}

func (r *authRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m model.UserModel
	query := `SELECT id, email, password_hash, first_name, last_name, role, manager_id, created_at, updated_at
	FROM users
	WHERE id=$1`
	if err := r.pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.Email, &m.PasswordHash, &m.FirstName, &m.LastName, &m.Role, &m.ManagerID, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("repository.GetByEmail(user): %w", err)
	}
	return converter.ToUserEntity(&m), nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, e domain.RefreshToken) error {
	m := converter.ToRefreshTokenModel(e)
	query := `INSERT INTO refresh_tokens(user_id, jti, token_hash, expires_at, revoked, created_at) VALUES($1,$2,$3,$4,$5,$6)`
	if _, err := r.pool.Exec(ctx, query, m.UserID, m.JTI, m.TokenHash, m.ExpiresAt, m.Revoked, m.CreatedAt); err != nil {
		return fmt.Errorf("repository.SaveRefreshToken: %w", err)
	}
	return nil
}
