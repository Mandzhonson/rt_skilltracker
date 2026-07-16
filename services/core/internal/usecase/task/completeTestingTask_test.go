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

type mockPlanCompletionService struct {
	generateSkillsIfCompletedFunc func(ctx context.Context, planID uuid.UUID) error
}

func (m *mockPlanCompletionService) GenerateSkillsIfCompleted(ctx context.Context, planID uuid.UUID) error {
	if m.generateSkillsIfCompletedFunc != nil {
		return m.generateSkillsIfCompletedFunc(ctx, planID)
	}
	return nil
}

func TestTaskService_CompleteTestingTask(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		taskRepo *mock_postgres.MockTaskRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	otherEmployeeID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		planID uuid.UUID
		userID uuid.UUID

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:   "Успешное завершение тестовой задачи",
			planID: planID,
			userID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					CompleteTestingTask(gomock.Any(), planID).
					Return(nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(100, nil)
			},

			expectedErr: nil,
		},
		{
			name:   "План не найден",
			planID: planID,
			userID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedErr: postgres.ErrPlanNotFound,
		},
		{
			name:   "План принадлежит другому сотруднику - запрещено",
			planID: planID,
			userID: otherEmployeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedErr: ErrEmployeeForbidden,
		},
		{
			name:   "Ошибка при завершении тестовой задачи",
			planID: planID,
			userID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					CompleteTestingTask(gomock.Any(), planID).
					Return(errors.New("complete task error"))
			},

			expectedErr: errors.New("complete task error"),
		},
		{
			name:   "Ошибка при пересчете прогресса",
			planID: planID,
			userID: employeeID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: employeeID,
					CreatedBy:  managerID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					CompleteTestingTask(gomock.Any(), planID).
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

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)

			testCase.mockBehavior(planRepo, taskRepo)

			mockPlanCompletion := &mockPlanCompletionService{}

			src := NewTaskService(
				taskRepo,
				planRepo,
				mockPlanCompletion,
			)

			err := src.CompleteTestingTask(
				context.Background(),
				testCase.planID,
				testCase.userID,
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
