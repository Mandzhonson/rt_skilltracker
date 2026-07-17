package admin

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAdminService_AdminGetPlan(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	managerID := uuid.New()
	task1ID := uuid.New()
	task2ID := uuid.New()

	now := time.Now()

	testTable := []struct {
		name string

		planID uuid.UUID

		mockBehavior mockBehavior

		expectedPlan *domain.PlanWithTasks
		expectedErr  error
	}{
		{
			name:   "Успешное получение плана",
			planID: planID,

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
					GetByIDWithTasks(gomock.Any(), planID).
					Return(planWithTasks, nil)
			},

			expectedPlan: &domain.PlanWithTasks{
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
			name:   "План не найден",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByIDWithTasks(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedPlan: nil,
			expectedErr:  postgres.ErrPlanNotFound,
		},
		{
			name:   "Ошибка при получении плана",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByIDWithTasks(gomock.Any(), planID).
					Return(nil, errors.New("database error"))
			},

			expectedPlan: nil,
			expectedErr:  errors.New("database error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)

			testCase.mockBehavior(planRepo)

			// Создаем сервис с правильным количеством аргументов
			src := NewAdminService(
				nil,      // userRepo
				planRepo, // planRepo
				nil,      // skillRepo
				nil,      // storage
			)

			result, err := src.AdminGetPlan(
				context.Background(),
				testCase.planID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testCase.expectedPlan.Plan.ID, result.Plan.ID)
				assert.Equal(t, testCase.expectedPlan.Plan.Title, result.Plan.Title)
				assert.Equal(t, len(testCase.expectedPlan.Tasks), len(result.Tasks))
				if len(result.Tasks) > 0 {
					assert.Equal(t, testCase.expectedPlan.Tasks[0].Title, result.Tasks[0].Title)
					assert.Equal(t, testCase.expectedPlan.Tasks[0].Status, result.Tasks[0].Status)
				}
			}
		})
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
