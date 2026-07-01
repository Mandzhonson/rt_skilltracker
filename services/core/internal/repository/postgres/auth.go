package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	Create(ctx context.Context, u *domain.User) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

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
	query := `SELECT id, email, password_hash, first_name, last_name, role, manager_id, created_at, updated_at
	FROM users
	WHERE email=$1`
	if err := r.pool.QueryRow(ctx, query, email).Scan(&m.ID, &m.Email, &m.PasswordHash, &m.FirstName, &m.LastName, &m.Role, &m.ManagerID, &m.CreatedAt, &m.UpdatedAt); err != nil {
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
		return nil, fmt.Errorf("repository.GetByEmail(user): %w", err)
	}
	return converter.ToUserEntity(&m), nil
}
