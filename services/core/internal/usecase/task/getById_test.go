package task

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestTaskService_GetByID(t *testing.T) {
	type mockBehavior func(
		taskRepo *mock_postgres.MockTaskRepository,
		planRepo *mock_postgres.MockPlanRepository,
	)

	taskID := uuid.New()
	planID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()

	testTable := []struct {
		name string

		managerID uuid.UUID
		taskID    uuid.UUID

		mockBehavior mockBehavior

		expectedTask *domain.Task
		expectedErr  error
	}{
		{
			name:      "Успешное получение задачи",
			managerID: managerID,
			taskID:    taskID,

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:          taskID,
					PlanID:      planID,
					Title:       "Test Task",
					Description: strPtr("Test Description"),
					Position:    1,
					Status:      domain.TaskTodo,
				}

				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedTask: &domain.Task{
				ID:          taskID,
				PlanID:      planID,
				Title:       "Test Task",
				Description: strPtr("Test Description"),
				Position:    1,
				Status:      domain.TaskTodo,
			},
			expectedErr: nil,
		},
		{
			name:      "Неверный ID задачи (пустой UUID)",
			managerID: managerID,
			taskID:    uuid.Nil,

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedTask: nil,
			expectedErr:  ErrInvalidTaskID,
		},
		{
			name:      "Задача не найдена",
			managerID: managerID,
			taskID:    taskID,

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(nil, postgres.ErrTaskNotFound)
			},

			expectedTask: nil,
			expectedErr:  ErrTaskNotFound,
		},
		{
			name:      "План не найден",
			managerID: managerID,
			taskID:    taskID,

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Test Task",
					Status: domain.TaskTodo,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedTask: nil,
			expectedErr:  postgres.ErrPlanNotFound,
		},
		{
			name:      "Нет прав (другой менеджер)",
			managerID: otherManagerID,
			taskID:    taskID,

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Test Task",
					Status: domain.TaskTodo,
				}

				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedTask: nil,
			expectedErr:  ErrManagerForbidden,
		},
		{
			name:      "Ошибка при получении задачи",
			managerID: managerID,
			taskID:    taskID,

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(nil, assert.AnError)
			},

			expectedTask: nil,
			expectedErr:  assert.AnError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)

			testCase.mockBehavior(taskRepo, planRepo)

			mockPlanCompletion := &mockPlanCompletionService{}

			src := NewTaskService(
				taskRepo,
				planRepo,
				mockPlanCompletion,
			)

			task, err := src.GetByID(
				context.Background(),
				testCase.managerID,
				testCase.taskID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, testCase.expectedTask.ID, task.ID)
				assert.Equal(t, testCase.expectedTask.Title, task.Title)
				assert.Equal(t, testCase.expectedTask.PlanID, task.PlanID)
				assert.Equal(t, testCase.expectedTask.Status, task.Status)
			}
		})
	}
}
