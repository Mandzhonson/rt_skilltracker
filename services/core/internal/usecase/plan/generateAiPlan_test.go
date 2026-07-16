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

type MockAIClientForGenerate struct {
	generateFunc     func(ctx context.Context, prompt string) (string, error)
	generateTestFunc func(ctx context.Context, prompt string) (string, error)
}

func (m *MockAIClientForGenerate) Generate(ctx context.Context, prompt string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, prompt)
	}
	return `{
		"title": "Generated Plan",
		"description": "Generated Description",
		"tasks": [
			{"title": "Task 1", "description": "Description 1"},
			{"title": "Task 2", "description": "Description 2"}
		]
	}`, nil
}

func (m *MockAIClientForGenerate) GenerateTest(ctx context.Context, prompt string) (string, error) {
	if m.generateTestFunc != nil {
		return m.generateTestFunc(ctx, prompt)
	}
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

func TestPlanService_generateAIPlan(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		userRepo *mock_postgres.MockUserRepository,
		taskRepo *mock_postgres.MockTaskRepository,
		skillRepo *mock_postgres.MockSkillRepository,
		testRepo *mock_postgres.MockTestRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		planID uuid.UUID

		mockBehavior mockBehavior

		expectError bool
	}{
		{
			name:   "Успешная генерация AI плана",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      strPtr("Test Description"),
					GenerationStatus: domain.GenerationPending,
					CreationType:     domain.CreationAI,
					Progress:         0,
					Status:           domain.PlanActive,
				}

				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return([]*domain.Skill{}, nil)

				planRepo.EXPECT().
					UpdateAIContent(gomock.Any(), planID, "Generated Plan", gomock.Any()).
					Return(nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil).
					Times(3)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil).
					AnyTimes()

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil).
					AnyTimes()

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return([]*domain.Skill{}, nil).
					AnyTimes()

				testRepo.EXPECT().
					CreateWithQuestions(gomock.Any(), gomock.Any(), gomock.Any()).
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
			name:   "План не найден",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectError: true,
		},
		{
			name:   "Сотрудник не найден - статус failed",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(nil, postgres.ErrUserNotFound)

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil)
			},

			expectError: true,
		},
		{
			name:   "Ошибка при получении навыков - статус failed",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				employee := &domain.User{
					ID:       employeeID,
					Email:    "employee@mail.ru",
					Role:     domain.RoleEmployee,
					Position: "Developer",
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil)
			},

			expectError: true,
		},
		{
			name:   "Ошибка при генерации плана AI - статус failed",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				employee := &domain.User{
					ID:       employeeID,
					Email:    "employee@mail.ru",
					Role:     domain.RoleEmployee,
					Position: "Developer",
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return([]*domain.Skill{}, nil)

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), planID, domain.GenerationFailed).
					Return(nil)
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
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)
			testRepo := mock_postgres.NewMockTestRepository(ctrl)

			testCase.mockBehavior(planRepo, userRepo, taskRepo, skillRepo, testRepo)

			mockClient := &MockAIClientForGenerate{
				generateFunc: func(ctx context.Context, prompt string) (string, error) {
					if testCase.name == "Ошибка при генерации плана AI - статус failed" {
						return "", errors.New("AI generation error")
					}
					return `{
						"title": "Generated Plan",
						"description": "Generated Description",
						"tasks": [
							{"title": "Task 1", "description": "Description 1"},
							{"title": "Task 2", "description": "Description 2"}
						]
					}`, nil
				},
				generateTestFunc: func(ctx context.Context, prompt string) (string, error) {
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
				},
			}
			aiService := ai.NewAiService(mockClient)

			src := &planService{
				planRepo:  planRepo,
				userRepo:  userRepo,
				taskRepo:  taskRepo,
				skillRepo: skillRepo,
				testRepo:  testRepo,
				aiService: aiService,
			}

			src.generateAIPlan(
				context.Background(),
				testCase.planID,
			)

		})
	}
}
