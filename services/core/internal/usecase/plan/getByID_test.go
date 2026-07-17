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

func TestPlanService_GetByID(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
	)

	planID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	employeeID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		managerID uuid.UUID
		planID    uuid.UUID

		mockBehavior mockBehavior

		expectedPlan *domain.Plan
		expectedErr  error
	}{
		{
			name:      "Успешное получение плана",
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
			},

			expectedPlan: &domain.Plan{
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
			expectedErr: nil,
		},
		{
			name:      "Неверный ID плана (пустой UUID)",
			managerID: managerID,
			planID:    uuid.Nil,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedPlan: nil,
			expectedErr:  ErrInvalidPlanID,
		},
		{
			name:      "План не найден",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedPlan: nil,
			expectedErr:  ErrPlanNotFound,
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

			expectedPlan: nil,
			expectedErr:  ErrForbidden,
		},
		{
			name:      "Ошибка при получении плана из БД",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, assert.AnError)
			},

			expectedPlan: nil,
			expectedErr:  assert.AnError,
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

			plan, err := src.GetByID(
				context.Background(),
				testCase.managerID,
				testCase.planID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, plan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, plan)
				assert.Equal(t, testCase.expectedPlan.ID, plan.ID)
				assert.Equal(t, testCase.expectedPlan.Title, plan.Title)
				assert.Equal(t, testCase.expectedPlan.CreatedBy, plan.CreatedBy)
				assert.Equal(t, testCase.expectedPlan.Status, plan.Status)
			}
		})
	}
}
