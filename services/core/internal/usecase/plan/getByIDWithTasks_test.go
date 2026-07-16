package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlanService_GetByIDWithTasks(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
	)

	planID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	employeeID := uuid.New()
	task1ID := uuid.New()
	task2ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		managerID uuid.UUID
		planID    uuid.UUID

		mockBehavior mockBehavior

		expectedPlanWithTasks *domain.PlanWithTasks
		expectedErr           error
		expectedErrContains   string
	}{
		{
			name:      "Успешное получение плана с задачами",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

				planWithTasks := &domain.PlanWithTasks{
					Plan:  plan,
					Tasks: tasks,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				planRepo.EXPECT().
					GetByIDWithTasks(gomock.Any(), planID).
					Return(planWithTasks, nil)
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
			name:      "План не найден при GetByID",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedPlanWithTasks: nil,
			expectedErr:           ErrPlanNotFound,
		},
		{
			name:      "Нет прав (другой менеджер)",
			managerID: otherManagerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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
			expectedErr:           ErrManagerForbidden,
		},
		{
			name:      "План не найден при GetByIDWithTasks",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

				planRepo.EXPECT().
					GetByIDWithTasks(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedPlanWithTasks: nil,
			expectedErr:           ErrPlanNotFound,
		},
		{
			name:      "Ошибка при GetByID",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, errors.New("database error"))
			},

			expectedPlanWithTasks: nil,
			expectedErrContains:   "get plan by id: database error",
		},
		{
			name:      "Ошибка при GetByIDWithTasks",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

				planRepo.EXPECT().
					GetByIDWithTasks(gomock.Any(), planID).
					Return(nil, errors.New("database error"))
			},

			expectedPlanWithTasks: nil,
			expectedErrContains:   "get plan with tasks: database error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)

			testCase.mockBehavior(planRepo)

			src := NewPlanService(
				planRepo,
				nil, // userRepo
				nil, // taskRepo
				nil, // skillRepo
				nil, // testRepo
				nil, // aiService
			)

			result, err := src.GetByIDWithTasks(
				context.Background(),
				testCase.managerID,
				testCase.planID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, result)
			} else if testCase.expectedErrContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedErrContains)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testCase.expectedPlanWithTasks.Plan.ID, result.Plan.ID)
				assert.Equal(t, testCase.expectedPlanWithTasks.Plan.Title, result.Plan.Title)
				assert.Equal(t, len(testCase.expectedPlanWithTasks.Tasks), len(result.Tasks))
				if len(result.Tasks) > 0 {
					assert.Equal(t, testCase.expectedPlanWithTasks.Tasks[0].Title, result.Tasks[0].Title)
					assert.Equal(t, testCase.expectedPlanWithTasks.Tasks[0].Status, result.Tasks[0].Status)
				}
			}
		})
	}
}
