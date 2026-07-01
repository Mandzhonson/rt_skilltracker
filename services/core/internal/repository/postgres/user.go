package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (uuid.UUID, error)
	// GetById(ctx context.Context, id string) (*domain.User, error)
	// GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	m := converter.ToUserModel(u)
	query := `INSERT INTO users(email, password_hash, first_name, last_name, manager_id) VALUES($1,$2,$3,$4,$5) RETURNING id`
	err := r.pool.QueryRow(ctx, query, m.Email, m.PasswordHash, m.FirstName, m.LastName, m.ManagerID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repository.Create (user): %w", err)
	}
	return id, nil
}
