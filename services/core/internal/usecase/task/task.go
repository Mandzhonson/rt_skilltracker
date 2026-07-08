package task

import (
	"core_service/internal/repository/postgres"
	"errors"
)

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrInvalidTaskID = errors.New("invalid task id")
	ErrInvalidTitle  = errors.New("invalid title")
	ErrPlanNotFound  = errors.New("plan not found")
	ErrForbidden     = errors.New("forbidden")
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
