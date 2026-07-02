package postgres

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Create(ctx context.Context, u *domain.User) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	SaveRefreshToken(ctx context.Context, e domain.RefreshToken) error
}
