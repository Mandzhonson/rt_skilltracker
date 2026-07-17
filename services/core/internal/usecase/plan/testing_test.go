package plan

import (
	"context"
	"core_service/internal/domain"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/ai"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockAIClientForAttach - мок для AI клиента
type mockAIClientForAttach struct{}

func (m *mockAIClientForAttach) Generate(ctx context.Context, prompt string) (string, error) {
	// Возвращаем 10 вопросов (как ожидает GenerateTest)
	return `{
		"questions": [
			{"question":"Q1","option_a":"A1","option_b":"B1","option_c":"C1","option_d":"D1","correct_option":"A"},
			{"question":"Q2","option_a":"A2","option_b":"B2","option_c":"C2","option_d":"D2","correct_option":"B"},
			{"question":"Q3","option_a":"A3","option_b":"B3","option_c":"C3","option_d":"D3","correct_option":"C"},
			{"question":"Q4","option_a":"A4","option_b":"B4","option_c":"C4","option_d":"D4","correct_option":"D"},
			{"question":"Q5","option_a":"A5","option_b":"B5","option_c":"C5","option_d":"D5","correct_option":"A"},
			{"question":"Q6","option_a":"A6","option_b":"B6","option_c":"C6","option_d":"D6","correct_option":"B"},
			{"question":"Q7","option_a":"A7","option_b":"B7","option_c":"C7","option_d":"D7","correct_option":"C"},
			{"question":"Q8","option_a":"A8","option_b":"B8","option_c":"C8","option_d":"D8","correct_option":"D"},
			{"question":"Q9","option_a":"A9","option_b":"B9","option_c":"C9","option_d":"D9","correct_option":"A"},
			{"question":"Q10","option_a":"A10","option_b":"B10","option_c":"C10","option_d":"D10","correct_option":"B"}
		]
	}`, nil
}

func TestPlanService_attachTesting(t *testing.T) {
	type mockBehavior func(
		taskRepo *mock_postgres.MockTaskRepository,
		testRepo *mock_postgres.MockTestRepository,
		userRepo *mock_postgres.MockUserRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	managerID := uuid.New()
	testID := uuid.New()
	description := "Test Description"

	testTable := []struct {
		name string

		plan *domain.Plan

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешное добавление тестовой задачи",
			plan: &domain.Plan{
				ID:          planID,
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Title:       "Test Plan",
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				userRepo *mock_postgres.MockUserRepository,
			) {
				existingTasks := []*domain.Task{
					{
						ID:          uuid.New(),
						PlanID:      planID,
						Title:       "Task 1",
						Description: strPtr("Task 1 Description"),
						Position:    1,
						Status:      domain.TaskDone,
					},
				}

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(existingTasks, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil)

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
					Return(employee, nil)

				testRepo.EXPECT().
					CreateWithQuestions(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(testID, nil)
			},

			expectedErr: nil,
		},
		{
			name: "Ошибка при получении списка задач",
			plan: &domain.Plan{
				ID:          planID,
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Title:       "Test Plan",
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				userRepo *mock_postgres.MockUserRepository,
			) {
				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(nil, errors.New("database error"))
			},

			expectedErr: errors.New("database error"),
		},
		{
			name: "Ошибка при создании тестовой задачи",
			plan: &domain.Plan{
				ID:          planID,
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Title:       "Test Plan",
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				userRepo *mock_postgres.MockUserRepository,
			) {
				existingTasks := []*domain.Task{
					{
						ID:          uuid.New(),
						PlanID:      planID,
						Title:       "Task 1",
						Description: strPtr("Task 1 Description"),
						Position:    1,
						Status:      domain.TaskDone,
					},
				}

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(existingTasks, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("create task error"))
			},

			expectedErr: errors.New("create task error"),
		},
		{
			name: "Ошибка при генерации теста",
			plan: &domain.Plan{
				ID:          planID,
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Title:       "Test Plan",
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				testRepo *mock_postgres.MockTestRepository,
				userRepo *mock_postgres.MockUserRepository,
			) {
				existingTasks := []*domain.Task{
					{
						ID:          uuid.New(),
						PlanID:      planID,
						Title:       "Task 1",
						Description: strPtr("Task 1 Description"),
						Position:    1,
						Status:      domain.TaskDone,
					},
				}

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(existingTasks, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil)

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
					Return(employee, nil)

				testRepo.EXPECT().
					CreateWithQuestions(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("test creation error"))
			},

			expectedErr: errors.New("test creation error"),
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

			testCase.mockBehavior(taskRepo, testRepo, userRepo)

			mockClient := &mockAIClientForAttach{}
			aiService := ai.NewAiService(mockClient)

			src := &planService{
				planRepo:  planRepo,
				userRepo:  userRepo,
				taskRepo:  taskRepo,
				skillRepo: skillRepo,
				testRepo:  testRepo,
				aiService: aiService,
			}

			err := src.attachTesting(
				context.Background(),
				testCase.plan,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
