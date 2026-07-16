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

func TestTaskService_UpdateStatus(t *testing.T) {
	type mockBehavior func(
		taskRepo *mock_postgres.MockTaskRepository,
		planRepo *mock_postgres.MockPlanRepository,
	)

	taskID := uuid.New()
	planID := uuid.New()
	employeeID := uuid.New()
	otherEmployeeID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		input UpdateTaskStatusInput

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешное обновление статуса задачи (не завершена)",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskInProgress,
			},

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
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					UpdateStatus(gomock.Any(), taskID, string(domain.TaskInProgress)).
					Return(nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(50, nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное обновление статуса задачи (завершена, генерируем навыки)",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskDone,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Test Task",
					Status: domain.TaskInProgress,
				}

				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					UpdateStatus(gomock.Any(), taskID, string(domain.TaskDone)).
					Return(nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(100, nil)
			},

			expectedErr: nil,
		},
		{
			name: "Неверный ID задачи (пустой UUID)",
			input: UpdateTaskStatusInput{
				TaskID: uuid.Nil,
				UserID: employeeID,
				Status: domain.TaskInProgress,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidTaskID,
		},
		{
			name: "Неверный ID пользователя (пустой UUID)",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: uuid.Nil,
				Status: domain.TaskInProgress,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidUserID,
		},
		{
			name: "Неверный статус",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: "invalid_status",
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidStatus,
		},
		{
			name: "Задача не найдена",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskInProgress,
			},

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
			name: "Статус задачи уже такой же",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskTodo,
			},

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
			},

			expectedErr: ErrInvalidStatus,
		},
		{
			name: "План не найден",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskInProgress,
			},

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

			expectedErr: ErrPlanNotFound,
		},
		{
			name: "Нет прав (план принадлежит другому сотруднику)",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: otherEmployeeID,
				Status: domain.TaskInProgress,
			},

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
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedErr: ErrEmployeeForbidden,
		},
		{
			name: "Ошибка при обновлении статуса задачи",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskInProgress,
			},

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
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					UpdateStatus(gomock.Any(), taskID, string(domain.TaskInProgress)).
					Return(errors.New("update status error"))
			},

			expectedErr: errors.New("update status error"),
		},
		{
			name: "Ошибка при пересчете прогресса",
			input: UpdateTaskStatusInput{
				TaskID: taskID,
				UserID: employeeID,
				Status: domain.TaskInProgress,
			},

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
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				taskRepo.EXPECT().
					GetByID(gomock.Any(), taskID).
					Return(task, nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					UpdateStatus(gomock.Any(), taskID, string(domain.TaskInProgress)).
					Return(nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(0, errors.New("recalculate error"))
			},

			expectedErr: errors.New("recalculate error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)

			testCase.mockBehavior(taskRepo, planRepo)

			mockPlanCompletion := &mockPlanCompletionService{
				generateSkillsIfCompletedFunc: func(ctx context.Context, planID uuid.UUID) error {
					return nil
				},
			}

			src := NewTaskService(
				taskRepo,
				planRepo,
				mockPlanCompletion,
			)

			err := src.UpdateStatus(
				context.Background(),
				testCase.input,
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
