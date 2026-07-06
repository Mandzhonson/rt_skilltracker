package postgres

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, e domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, jti string) error
}

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, update *domain.UpdateUserProfile) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashPassword string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarKey *string) error
}
