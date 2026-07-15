package task

import (
	"context"
	"core_service/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
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

type PlanCompletionService interface {
	GenerateSkillsIfCompleted(ctx context.Context, planID uuid.UUID) error
}

type taskService struct {
	taskRepo              postgres.TaskRepository
	planRepo              postgres.PlanRepository
	planCompletionService PlanCompletionService
}

func NewTaskService(taskRepo postgres.TaskRepository, planRepo postgres.PlanRepository, planCompletionService PlanCompletionService) *taskService {
	return &taskService{
		taskRepo:              taskRepo,
		planRepo:              planRepo,
		planCompletionService: planCompletionService,
	}
}
