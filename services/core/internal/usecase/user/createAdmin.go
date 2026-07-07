package user

import (
	"context"
	"core_service/internal/domain"
)

func (s *userService) CreateAdmin(ctx context.Context, input CreateUserInput) error {

	exists, err := s.userRepo.ExistsAdmin(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	admin := domain.NewAdmin(
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
	)

	_, err = s.CreateUser(
		ctx,
		CreateUserInput{
			Email:     admin.Email,
			Password:  input.Password,
			FirstName: admin.FirstName,
			LastName:  admin.LastName,
			Role:      domain.RoleAdmin,
		},
	)

	return err
}
