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

func TestTaskService_Create(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		taskRepo *mock_postgres.MockTaskRepository,
	)

	planID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	taskID := uuid.New()
	title := "New Task"
	description := "Task Description"

	testTable := []struct {
		name string

		input CreateTaskInput

		mockBehavior mockBehavior

		expectedID       uuid.UUID
		expectedErr      error
		expectedContains string
	}{
		{
			name: "Успешное создание задачи",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(1, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(taskID, nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(50, nil)
			},

			expectedID:  taskID,
			expectedErr: nil,
		},
		{
			name: "Успешное создание задачи (без описания)",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: nil,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(2, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(taskID, nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(50, nil)
			},

			expectedID:  taskID,
			expectedErr: nil,
		},
		{
			name: "Пустое название задачи",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       "   ",
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidTitle,
		},
		{
			name: "План не найден",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrPlanNotFound,
		},
		{
			name: "Нет прав (другой менеджер)",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   otherManagerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrManagerForbidden,
		},
		{
			name: "План архивирован",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanArchived,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrPlanArchived,
		},
		{
			name: "Ошибка при получении позиции",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(0, errors.New("position error"))
			},

			expectedID:  uuid.Nil,
			expectedErr: errors.New("position error"),
		},
		{
			name: "Ошибка при создании задачи",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(1, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("create error"))
			},

			expectedID:       uuid.Nil,
			expectedContains: "create task: create error",
		},
		{
			name: "Ошибка при пересчете прогресса",
			input: CreateTaskInput{
				PlanID:      planID,
				CreatedBy:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Test Plan",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					GetNextPosition(gomock.Any(), planID).
					Return(1, nil)

				taskRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(taskID, nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(0, errors.New("recalculate error"))
			},

			expectedID:       uuid.Nil,
			expectedContains: "recalculate progress: recalculate error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)

			testCase.mockBehavior(planRepo, taskRepo)

			mockPlanCompletion := &mockPlanCompletionService{}

			src := NewTaskService(
				taskRepo,
				planRepo,
				mockPlanCompletion,
			)

			id, err := src.Create(
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
