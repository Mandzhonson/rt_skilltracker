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
	ErrEmployeeForbidden   = errors.New("employee has no access")
	ErrForbidden           = errors.New("forbidden")
	ErrManagerForbidden    = errors.New("manager has no access")
)

type planService struct {
	planRepo postgres.PlanRepository
	userRepo postgres.UserRepository
	taskRepo postgres.TaskRepository
}

func NewPlanService(planRepo postgres.PlanRepository, userRepo postgres.UserRepository, taskRepo postgres.TaskRepository) *planService {
	return &planService{
		planRepo: planRepo,
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}
