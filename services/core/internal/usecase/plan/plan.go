package plan

import (
	"core_service/internal/repository/postgres"
	"core_service/internal/usecase/ai"
	"errors"
	"strings"
)

var (
	ErrInvalidTitle         = errors.New("invalid title")
	ErrInvalidEmployee      = errors.New("invalid employee")
	ErrInvalidCreator       = errors.New("invalid creator")
	ErrEmployeeNotAssigned  = errors.New("employee is not assigned to manager")
	ErrPlanNotFound         = errors.New("plan not found")
	ErrInvalidPlanID        = errors.New("invalid plan id")
	ErrEmployeeForbidden    = errors.New("employee has no access")
	ErrForbidden            = errors.New("forbidden")
	ErrManagerForbidden     = errors.New("manager has no access")
	ErrTestGenerationFailed = errors.New("test generation failed")
)

type planService struct {
	planRepo  postgres.PlanRepository
	userRepo  postgres.UserRepository
	taskRepo  postgres.TaskRepository
	skillRepo postgres.SkillRepository
	testRepo  postgres.TestRepository
	aiService *ai.AiService
}

func NewPlanService(planRepo postgres.PlanRepository, userRepo postgres.UserRepository, taskRepo postgres.TaskRepository, skillRepo postgres.SkillRepository, testRepo postgres.TestRepository, aiService *ai.AiService) *planService {
	return &planService{
		planRepo:  planRepo,
		userRepo:  userRepo,
		taskRepo:  taskRepo,
		skillRepo: skillRepo,
		testRepo:  testRepo,
		aiService: aiService,
	}
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizeCorrectOption(option string) string {

	option = strings.ToUpper(
		strings.TrimSpace(option),
	)

	switch option {
	case "A", "B", "C", "D":
		return option
	default:
		return ""
	}
}
