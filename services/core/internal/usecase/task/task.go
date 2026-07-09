package task

import (
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"errors"
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTaskID     = errors.New("invalid task id")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidTitle      = errors.New("invalid title")
	ErrPlanNotFound      = errors.New("plan not found")
	ErrManagerForbidden  = errors.New("manager has no access")
	ErrEmployeeForbidden = errors.New("employee has no access")
	ErrInvalidUpdate     = errors.New("nothing to update")
	ErrInvalidStatus     = errors.New("invalid task status")
)

type taskService struct {
	taskRepo postgres.TaskRepository
	planRepo postgres.PlanRepository
}

func NewTaskService(taskRepo postgres.TaskRepository, planRepo postgres.PlanRepository) *taskService {
	return &taskService{
		taskRepo: taskRepo,
		planRepo: planRepo,
	}
}

func validateStatusChange(current domain.TaskStatus, next domain.TaskStatus) bool {
	switch current {
	case domain.TaskTodo:
		return next == domain.TaskInProgress
	case domain.TaskInProgress:
		return next == domain.TaskDone
	case domain.TaskDone:
		return false
	}
	return false
}
