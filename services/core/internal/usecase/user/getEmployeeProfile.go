package user

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

func (s *userService) GetEmployeeProfile(ctx context.Context, managerID uuid.UUID, employeeID uuid.UUID) (*EmployeeProfile, error) {

	manager, err := s.userRepo.GetById(ctx, managerID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !manager.IsManager() {
		return nil, ErrNotManager
	}

	employee, err := s.userRepo.GetById(ctx, employeeID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if employee.ManagerID == nil ||
		*employee.ManagerID != managerID {

		return nil, ErrForbidden
	}

	skills, err := s.skillRepo.ListByUserID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	plans, err := s.planRepo.ListAllByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	return &EmployeeProfile{
		User:   employee,
		Skills: skills,
		Plans:  plans,
	}, nil
}
