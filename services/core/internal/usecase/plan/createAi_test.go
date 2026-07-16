package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/ai"
	"core_service/internal/usecase/user"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockAIClient2 struct {
	generateFunc func(ctx context.Context, prompt string) (string, error)
}

func (m *MockAIClient2) Generate(ctx context.Context, prompt string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, prompt)
	}
	return `{
		"title": "Test AI Plan",
		"description": "Test Description",
		"tasks": [
			{"title": "Task 1", "description": "Description 1"},
			{"title": "Task 2", "description": "Description 2"}
		]
	}`, nil
}

func TestPlanService_CreateAI(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		userRepo *mock_postgres.MockUserRepository,
		taskRepo *mock_postgres.MockTaskRepository,
		skillRepo *mock_postgres.MockSkillRepository,
		testRepo *mock_postgres.MockTestRepository,
	)

	employeeID := uuid.New()
	managerID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		input CreateAIInput

		mockBehavior mockBehavior

		expectedID       uuid.UUID
		expectedErr      error
		expectedContains string
	}{
		{
			name: "Успешное создание AI плана",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "Learn Go from scratch",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				planRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(planID, nil)

				planRepo.EXPECT().
					UpdateGenerationStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					AnyTimes()

				planRepo.EXPECT().
					UpdateAIContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					AnyTimes()

				planRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&domain.Plan{
						ID:               planID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Go Programming",
						GenerationStatus: domain.GenerationPending,
						CreationType:     domain.CreationAI,
						Progress:         0,
						Status:           domain.PlanActive,
					}, nil).
					AnyTimes()

				userRepo.EXPECT().
					GetById(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:    employeeID,
						Email: "employee@mail.ru",
						Role:  domain.RoleEmployee,
					}, nil).
					AnyTimes()

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil).
					AnyTimes()

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), gomock.Any()).
					Return([]*domain.Skill{}, nil).
					AnyTimes()

				skillRepo.EXPECT().
					AttachToUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					AnyTimes()
			},

			expectedID:  planID,
			expectedErr: nil,
		},
		{
			name: "Пустая тема",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "   ",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidTitle,
		},
		{
			name: "Сотрудник не найден",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedID:  uuid.Nil,
			expectedErr: user.ErrUserNotFound,
		},
		{
			name: "Сотрудник не является сотрудником",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				admin := &domain.User{
					ID:    employeeID,
					Email: "admin@mail.ru",
					Role:  domain.RoleAdmin,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(admin, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidEmployee,
		},
		{
			name: "Создатель не найден",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				employee := &domain.User{
					ID:    employeeID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedID:  uuid.Nil,
			expectedErr: user.ErrUserNotFound,
		},
		{
			name: "Создатель не является менеджером",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				employee := &domain.User{
					ID:    employeeID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				notManager := &domain.User{
					ID:    managerID,
					Email: "notmanager@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(notManager, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidCreator,
		},
		{
			name: "Сотрудник не закреплен за менеджером",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					Role:      domain.RoleEmployee,
					ManagerID: nil,
				}

				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrEmployeeNotAssigned,
		},
		{
			name: "Сотрудник закреплен за другим менеджером",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				otherManagerID := uuid.New()
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					Role:      domain.RoleEmployee,
					ManagerID: &otherManagerID,
				}

				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrEmployeeNotAssigned,
		},
		{
			name: "Ошибка при создании плана в БД",
			input: CreateAIInput{
				EmployeeID:  employeeID,
				CreatedBy:   managerID,
				Topic:       "Go Programming",
				Description: "",
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				planRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("database error"))
			},

			expectedID:       uuid.Nil,
			expectedContains: "database error",
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

			mockClient := &MockAIClient2{
				generateFunc: func(ctx context.Context, prompt string) (string, error) {
					return `{
						"title": "Test AI Plan",
						"description": "Test Description",
						"tasks": [
							{"title": "Task 1", "description": "Description 1"},
							{"title": "Task 2", "description": "Description 2"}
						]
					}`, nil
				},
			}
			aiService := ai.NewAiService(mockClient)

			src := NewPlanService(
				planRepo,
				userRepo,
				taskRepo,
				skillRepo,
				testRepo,
				aiService,
			)

			id, err := src.CreateAI(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Equal(t, testCase.expectedID, id)
			} else if testCase.expectedContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedContains)
				assert.Equal(t, testCase.expectedID, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedID, id)
			}
		})
	}
}
