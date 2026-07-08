package plan

import (
	"core_service/internal/repository/postgres"
	"errors"
)

var (
	ErrInvalidTitle        = errors.New("invalid title")
	ErrInvalidEmployee     = errors.New("invalid employee")
	ErrInvalidCreator      = errors.New("invalid creator")
	ErrEmployeeNotAssigned = errors.New("employee is not assigned to manager")
	ErrPlanNotFound        = errors.New("plan not found")
	ErrInvalidPlanID       = errors.New("invalid plan id")
	ErrForbidden           = errors.New("forbidden")
)

type planService struct {
	planRepo postgres.PlanRepository
	userRepo postgres.UserRepository
}

func NewPlanService(planRepo postgres.PlanRepository, userRepo postgres.UserRepository) *planService {
	return &planService{
		planRepo: planRepo,
		userRepo: userRepo,
	}
}
