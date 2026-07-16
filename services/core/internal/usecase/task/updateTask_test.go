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

func TestTaskService_Update(t *testing.T) {
	type mockBehavior func(
		taskRepo *mock_postgres.MockTaskRepository,
		planRepo *mock_postgres.MockPlanRepository,
	)

	taskID := uuid.New()
	planID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	title := "Updated Task Title"
	description := "Updated Task Description"
	emptyTitle := ""

	testTable := []struct {
		name string

		input UpdateTaskInput

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешное обновление задачи (название и описание)",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       &title,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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
					Update(gomock.Any(), taskID, &title, &description).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное обновление задачи (только название)",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       &title,
				Description: nil,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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
					Update(gomock.Any(), taskID, &title, nil).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное обновление задачи (только описание)",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       nil,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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
					Update(gomock.Any(), taskID, nil, &description).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Неверный ID задачи (пустой UUID)",
			input: UpdateTaskInput{
				TaskID:      uuid.Nil,
				ManagerID:   managerID,
				Title:       &title,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidTaskID,
		},
		{
			name: "Нет данных для обновления",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       nil,
				Description: nil,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidUpdate,
		},
		{
			name: "Задача не найдена",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       &title,
				Description: &description,
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
			name: "План не найден",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       &title,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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
			name: "Нет прав (другой менеджер)",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   otherManagerID,
				Title:       &title,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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
			name: "Пустое название (только пробелы)",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       &emptyTitle,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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

			expectedErr: ErrInvalidTitle,
		},
		{
			name: "Ошибка при обновлении задачи",
			input: UpdateTaskInput{
				TaskID:      taskID,
				ManagerID:   managerID,
				Title:       &title,
				Description: &description,
			},

			mockBehavior: func(
				taskRepo *mock_postgres.MockTaskRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				task := &domain.Task{
					ID:     taskID,
					PlanID: planID,
					Title:  "Old Title",
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
					Update(gomock.Any(), taskID, &title, &description).
					Return(errors.New("update error"))
			},

			expectedErr: errors.New("update error"),
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

			err := src.Update(
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
