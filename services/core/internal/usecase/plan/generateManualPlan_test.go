package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/ai"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

type MockAIClientForManual struct{}

func (m *MockAIClientForManual) Generate(ctx context.Context, prompt string) (string, error) {
	return "", nil
}

func (m *MockAIClientForManual) GenerateTest(ctx context.Context, prompt string) (string, error) {
	return `{
		"questions": [
			{
				"question": "What is Go?",
				"option_a": "Programming language",
				"option_b": "Framework",
				"option_c": "Library",
				"option_d": "Tool",
				"correct_option": "A"
			}
		]
	}`, nil
}

func TestPlanService_generateManualPlan(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		userRepo *mock_postgres.MockUserRepository,
		taskRepo *mock_postgres.MockTaskRepository,
		testRepo *mock_postgres.MockTestRepository,
		skillRepo *mock_postgres.MockSkillRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		plan *domain.Plan

		mockBehavior mockBehavior

		expectError bool
	}{
		{
			name: "Успешная генерация manual плана с тестом",
			plan: &domain.Plan{
				ID:               planID,
				EmployeeID:       employeeID,
				CreatedBy:        managerID,
				Title:            "Test Plan",
				Description:      strPtr("Test Description"),
				GenerationStatus: domain.GenerationPending,
				CreationType:     domain.CreationManual,
				Progress:         0,
				Status:           domain.PlanActive,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil).
					AnyTimes()

				testRepo.EXPECT().
					CreateWithQuestions(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil).
					AnyTimes()

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(1, nil).
					AnyTimes()

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil).
					AnyTimes()

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationReady).
					Return(nil).
					AnyTimes()

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil).
					AnyTimes()
			},

			expectError: false,
		},
		{
			name: "Успешная генерация manual плана без теста",
			plan: &domain.Plan{
				ID:               planID,
				EmployeeID:       employeeID,
				CreatedBy:        managerID,
				Title:            "Test Plan",
				Description:      strPtr("Test Description"),
				GenerationStatus: domain.GenerationPending,
				CreationType:     domain.CreationManual,
				Progress:         0,
				Status:           domain.PlanActive,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil).
					AnyTimes()

				testRepo.EXPECT().
					CreateWithQuestions(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil).
					AnyTimes()

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(0, errors.New("no tasks")).
					AnyTimes()

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationReady).
					Return(nil).
					AnyTimes()

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil).
					AnyTimes()
			},

			expectError: false,
		},
		{
			name: "Ошибка при генерации теста - статус failed",
			plan: &domain.Plan{
				ID:               planID,
				EmployeeID:       employeeID,
				CreatedBy:        managerID,
				Title:            "Test Plan",
				Description:      strPtr("Test Description"),
				GenerationStatus: domain.GenerationPending,
				CreationType:     domain.CreationManual,
				Progress:         0,
				Status:           domain.PlanActive,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil).
					AnyTimes()

				testRepo.EXPECT().
					CreateWithQuestions(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("test creation failed")).
					AnyTimes()

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil).
					AnyTimes()
			},

			expectError: true,
		},
		{
			name: "Сотрудник не найден - статус failed",
			plan: &domain.Plan{
				ID:               planID,
				EmployeeID:       employeeID,
				CreatedBy:        managerID,
				Title:            "Test Plan",
				Description:      strPtr("Test Description"),
				GenerationStatus: domain.GenerationPending,
				CreationType:     domain.CreationManual,
				Progress:         0,
				Status:           domain.PlanActive,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(nil, postgres.ErrUserNotFound).
					AnyTimes()

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil).
					AnyTimes()
			},

			expectError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)
			testRepo := mock_postgres.NewMockTestRepository(ctrl)
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)

			testCase.mockBehavior(planRepo, userRepo, taskRepo, testRepo, skillRepo)

			mockClient := &MockAIClientForManual{}
			aiService := ai.NewAiService(mockClient)

			src := &planService{
				planRepo:  planRepo,
				userRepo:  userRepo,
				taskRepo:  taskRepo,
				skillRepo: skillRepo,
				testRepo:  testRepo,
				aiService: aiService,
			}

			src.generateManualPlan(
				context.Background(),
				testCase.plan,
			)

		})
	}
}
