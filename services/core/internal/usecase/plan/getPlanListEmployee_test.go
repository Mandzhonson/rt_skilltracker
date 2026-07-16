package plan

import (
	"context"
	"core_service/internal/domain"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlanService_ListEmployeePlans(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		taskRepo *mock_postgres.MockTaskRepository,
	)

	employeeID := uuid.New()
	managerID := uuid.New()
	plan1ID := uuid.New()
	plan2ID := uuid.New()
	task1ID := uuid.New()
	task2ID := uuid.New()
	task3ID := uuid.New()
	task4ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		employeeID uuid.UUID

		mockBehavior mockBehavior

		expectedPlans []*domain.PlanWithTasks
		expectedErr   error
	}{
		{
			name:       "Успешное получение списка планов сотрудника",
			employeeID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plans := []*domain.Plan{
					{
						ID:               plan1ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 1",
						Description:      strPtr("Plan 1 Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					{
						ID:               plan2ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 2",
						Description:      strPtr("Plan 2 Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         80,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
				}

				tasks1 := []*domain.Task{
					{
						ID:          task1ID,
						PlanID:      plan1ID,
						Title:       "Task 1",
						Description: strPtr("Task 1 Description"),
						Position:    1,
						Status:      domain.TaskTodo,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          task2ID,
						PlanID:      plan1ID,
						Title:       "Task 2",
						Description: strPtr("Task 2 Description"),
						Position:    2,
						Status:      domain.TaskInProgress,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}

				tasks2 := []*domain.Task{
					{
						ID:          task3ID,
						PlanID:      plan2ID,
						Title:       "Task 3",
						Description: strPtr("Task 3 Description"),
						Position:    1,
						Status:      domain.TaskDone,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          task4ID,
						PlanID:      plan2ID,
						Title:       "Task 4",
						Description: strPtr("Task 4 Description"),
						Position:    2,
						Status:      domain.TaskTodo,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}

				planRepo.EXPECT().
					ListByEmployeeID(gomock.Any(), employeeID).
					Return(plans, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), plan1ID).
					Return(tasks1, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), plan2ID).
					Return(tasks2, nil)
			},

			expectedPlans: []*domain.PlanWithTasks{
				{
					Plan: &domain.Plan{
						ID:               plan1ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 1",
						Description:      strPtr("Plan 1 Description"),
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
							PlanID:      plan1ID,
							Title:       "Task 1",
							Description: strPtr("Task 1 Description"),
							Position:    1,
							Status:      domain.TaskTodo,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
						{
							ID:          task2ID,
							PlanID:      plan1ID,
							Title:       "Task 2",
							Description: strPtr("Task 2 Description"),
							Position:    2,
							Status:      domain.TaskInProgress,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					},
				},
				{
					Plan: &domain.Plan{
						ID:               plan2ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 2",
						Description:      strPtr("Plan 2 Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         80,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					Tasks: []*domain.Task{
						{
							ID:          task3ID,
							PlanID:      plan2ID,
							Title:       "Task 3",
							Description: strPtr("Task 3 Description"),
							Position:    1,
							Status:      domain.TaskDone,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
						{
							ID:          task4ID,
							PlanID:      plan2ID,
							Title:       "Task 4",
							Description: strPtr("Task 4 Description"),
							Position:    2,
							Status:      domain.TaskTodo,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Пустой список планов",
			employeeID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				planRepo.EXPECT().
					ListByEmployeeID(gomock.Any(), employeeID).
					Return([]*domain.Plan{}, nil)
			},

			expectedPlans: []*domain.PlanWithTasks{},
			expectedErr:   nil,
		},
		{
			name:       "Ошибка при получении планов",
			employeeID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				planRepo.EXPECT().
					ListByEmployeeID(gomock.Any(), employeeID).
					Return(nil, assert.AnError)
			},

			expectedPlans: nil,
			expectedErr:   assert.AnError,
		},
		{
			name:       "Ошибка при получении задач для плана",
			employeeID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plans := []*domain.Plan{
					{
						ID:               plan1ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 1",
						Description:      strPtr("Plan 1 Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
				}

				planRepo.EXPECT().
					ListByEmployeeID(gomock.Any(), employeeID).
					Return(plans, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), plan1ID).
					Return(nil, assert.AnError)
			},

			expectedPlans: nil,
			expectedErr:   assert.AnError,
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

			result, err := src.ListEmployeePlans(
				context.Background(),
				testCase.employeeID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testCase.expectedPlans), len(result))
				if len(result) > 0 {
					assert.Equal(t, testCase.expectedPlans[0].Plan.ID, result[0].Plan.ID)
					assert.Equal(t, testCase.expectedPlans[0].Plan.Title, result[0].Plan.Title)
					assert.Equal(t, len(testCase.expectedPlans[0].Tasks), len(result[0].Tasks))
				}
			}
		})
	}
}
