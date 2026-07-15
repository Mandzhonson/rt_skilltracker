package admin

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

func (s *adminService) UpdatePosition(ctx context.Context, userID uuid.UUID, position string) error {
	if strings.TrimSpace(position) == "" {
		return ErrInvalidPosition
	}

	return s.userRepo.UpdatePosition(ctx, userID, position)
}
