package user

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"
	"io"

	"github.com/google/uuid"
)

func (s *userService) GetEmployeeAvatar(ctx context.Context, employeeID uuid.UUID, managerID uuid.UUID) (io.ReadCloser, string, error) {

	employee, err := s.userRepo.GetById(ctx, employeeID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	if employee.ManagerID == nil || *employee.ManagerID != managerID {
		return nil, "", ErrEmployeeNotAssigned
	}

	if employee.AvatarKey == nil {
		return nil, "", ErrAvatarNotFound
	}

	return s.storage.GetAvatar(ctx, *employee.AvatarKey)
}
