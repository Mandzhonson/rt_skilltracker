package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlanService_GetEmployeePlan(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		taskRepo *mock_postgres.MockTaskRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	otherEmployeeID := uuid.New()
	managerID := uuid.New()
	task1ID := uuid.New()
	task2ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		employeeID uuid.UUID
		planID     uuid.UUID

		mockBehavior mockBehavior

		expectedPlanWithTasks *domain.PlanWithTasks
		expectedErr           error
	}{
		{
			name:       "Успешное получение плана сотрудника",
			employeeID: employeeID,
			planID:     planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      strPtr("Test Description"),
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         50,
					Status:           domain.PlanActive,
					CreatedAt:        now,
					UpdatedAt:        now,
				}

				tasks := []*domain.Task{
					{
						ID:          task1ID,
						PlanID:      planID,
						Title:       "Task 1",
						Description: strPtr("Task 1 Description"),
						Position:    1,
						Status:      domain.TaskTodo,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          task2ID,
						PlanID:      planID,
						Title:       "Task 2",
						Description: strPtr("Task 2 Description"),
						Position:    2,
						Status:      domain.TaskInProgress,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(tasks, nil)
			},

			expectedPlanWithTasks: &domain.PlanWithTasks{
				Plan: &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      strPtr("Test Description"),
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         50,
					Status:           domain.PlanActive,
					CreatedAt:        now,
					UpdatedAt:        now,
				},
				Tasks: []*domain.Task{
					{
						ID:          task1ID,
						PlanID:      planID,
						Title:       "Task 1",
						Description: strPtr("Task 1 Description"),
						Position:    1,
						Status:      domain.TaskTodo,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          task2ID,
						PlanID:      planID,
						Title:       "Task 2",
						Description: strPtr("Task 2 Description"),
						Position:    2,
						Status:      domain.TaskInProgress,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Неверный ID плана (пустой UUID)",
			employeeID: employeeID,
			planID:     uuid.Nil,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
			},

			expectedPlanWithTasks: nil,
			expectedErr:           ErrInvalidPlanID,
		},
		{
			name:       "План не найден",
			employeeID: employeeID,
			planID:     planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedPlanWithTasks: nil,
			expectedErr:           postgres.ErrPlanNotFound,
		},
		{
			name:       "Нет прав (план принадлежит другому сотруднику)",
			employeeID: otherEmployeeID,
			planID:     planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      strPtr("Test Description"),
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         50,
					Status:           domain.PlanActive,
					CreatedAt:        now,
					UpdatedAt:        now,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedPlanWithTasks: nil,
			expectedErr:           ErrEmployeeForbidden,
		},
		{
			name:       "Ошибка при получении задач",
			employeeID: employeeID,
			planID:     planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      strPtr("Test Description"),
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         50,
					Status:           domain.PlanActive,
					CreatedAt:        now,
					UpdatedAt:        now,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(nil, assert.AnError)
			},

			expectedPlanWithTasks: nil,
			expectedErr:           assert.AnError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)

			testCase.mockBehavior(planRepo, taskRepo)

			src := NewPlanService(
				planRepo,
				nil, // userRepo
				taskRepo,
				nil, // skillRepo
				nil, // testRepo
				nil, // aiService
			)

			result, err := src.GetEmployeePlan(
				context.Background(),
				testCase.employeeID,
				testCase.planID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testCase.expectedPlanWithTasks.Plan.ID, result.Plan.ID)
				assert.Equal(t, testCase.expectedPlanWithTasks.Plan.Title, result.Plan.Title)
				assert.Equal(t, testCase.expectedPlanWithTasks.Plan.EmployeeID, result.Plan.EmployeeID)
				assert.Equal(t, len(testCase.expectedPlanWithTasks.Tasks), len(result.Tasks))
				if len(result.Tasks) > 0 {
					assert.Equal(t, testCase.expectedPlanWithTasks.Tasks[0].Title, result.Tasks[0].Title)
				}
			}
		})
	}
}
