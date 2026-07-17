package admin

import (
	"context"

	"github.com/google/uuid"

	"core_service/internal/domain"
)

func (s *adminService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetById(ctx, id)
}
