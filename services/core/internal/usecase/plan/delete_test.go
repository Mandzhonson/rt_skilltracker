package plan

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

func TestPlanService_Delete(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
	)

	managerID := uuid.New()
	otherManagerID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		managerID uuid.UUID
		planID    uuid.UUID

		mockBehavior mockBehavior

		expectedErr      error
		expectedContains string
	}{
		{
			name:      "Успешное удаление плана",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

				planRepo.EXPECT().
					Delete(gomock.Any(), planID).
					Return(nil)
			},

			expectedErr: nil,
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

			expectedErr: ErrPlanNotFound,
		},
		{
			name:      "Нет прав на удаление (другой менеджер)",
			managerID: otherManagerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

			expectedErr: ErrManagerForbidden,
		},
		{
			name:      "Ошибка при удалении плана из БД",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

				planRepo.EXPECT().
					Delete(gomock.Any(), planID).
					Return(assert.AnError)
			},

			expectedContains: "delete plan:",
		},
		{
			name:      "Ошибка при получении плана (не ErrPlanNotFound)",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, assert.AnError)
			},

			expectedContains: "get plan:",
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

			err := src.Delete(
				context.Background(),
				testCase.managerID,
				testCase.planID,
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