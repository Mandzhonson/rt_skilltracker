package task

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_Delete(t *testing.T) {
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

		expectedErr      error
		expectedContains string
	}{
		{
			name:      "Успешное удаление задачи",
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

				taskRepo.EXPECT().
					Delete(gomock.Any(), taskID).
					Return(nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(50, nil)
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

			expectedErr: ErrInvalidTaskID,
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

			expectedErr: ErrTaskNotFound,
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

			expectedErr: postgres.ErrPlanNotFound,
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

			expectedErr: ErrManagerForbidden,
		},
		{
			name:      "Ошибка при удалении задачи",
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

				taskRepo.EXPECT().
					Delete(gomock.Any(), taskID).
					Return(errors.New("delete error"))
			},

			expectedErr: errors.New("delete error"),
		},
		{
			name:      "Ошибка при пересчете прогресса",
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

				taskRepo.EXPECT().
					Delete(gomock.Any(), taskID).
					Return(nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(0, errors.New("recalculate error"))
			},

			expectedContains: "recalculate progress: recalculate error",
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

			err := src.Delete(
				context.Background(),
				testCase.managerID,
				testCase.taskID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else if testCase.expectedContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
