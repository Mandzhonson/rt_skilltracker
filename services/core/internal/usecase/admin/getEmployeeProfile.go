package admin

import (
	"context"
	"errors"
	"fmt"

	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"core_service/internal/usecase/user"

	"github.com/google/uuid"
)

func (s *adminService) GetEmployeeProfileForAdmin(
	ctx context.Context,
	employeeID uuid.UUID,
) (*domain.EmployeeProfile, error) {

	employee, err := s.userRepo.GetById(ctx, employeeID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if !employee.IsEmployee() {
		return nil, ErrInvalidEmployee
	}

	plans, err := s.planRepo.ListAllByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}

	skills, err := s.skillRepo.ListByUserID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
	}

	return &domain.EmployeeProfile{
		User:   employee,
		Plans:  plans,
		Skills: skills,
	}, nil
}
