package test

import (
	"core_service/internal/repository/postgres"
	"core_service/internal/usecase/task"
	"errors"
)

var (
	ErrTestNotFound      = errors.New("test not found")
	ErrInvalidPlanID     = errors.New("invalid plan id")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidAnswers    = errors.New("invalid answers")
	ErrTestAlreadyPassed = errors.New("test already passed")
)

type testService struct {
	testRepo              postgres.TestRepository
	taskService           task.TaskService
	planRepo              postgres.PlanRepository
	planCompletionService task.PlanCompletionService
}

func NewTestService(testRepo postgres.TestRepository, taskService task.TaskService, planRepo postgres.PlanRepository, planCompletionService task.PlanCompletionService) *testService {
	return &testService{
		testRepo:              testRepo,
		taskService:           taskService,
		planRepo:              planRepo,
		planCompletionService: planCompletionService,
	}
}
