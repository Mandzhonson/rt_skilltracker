package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/ai"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlanService_Archive(t *testing.T) {
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

		expectedErr error
	}{
		{
			name:      "Успешное архивирование плана",
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
					Archive(gomock.Any(), planID).
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

			expectedErr: postgres.ErrPlanNotFound,
		},
		{
			name:      "Нет прав на архивирование (другой менеджер)",
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

			expectedErr: ErrForbidden,
		},
		{
			name:      "Ошибка при архивировании",
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
					Archive(gomock.Any(), planID).
					Return(assert.AnError)
			},

			expectedErr: assert.AnError,
		},
		{
			name:      "Архивирование уже архивированного плана",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
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

				planRepo.EXPECT().
					Archive(gomock.Any(), planID).
					Return(nil)
			},

			expectedErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)

			testCase.mockBehavior(planRepo)

			mockClient := &MockAIClient{}
			aiService := ai.NewAiService(mockClient)

			src := NewPlanService(
				planRepo,  // PlanRepository
				nil,       // UserRepository
				nil,       // TaskRepository
				nil,       // SkillRepository
				nil,       // TestRepository
				aiService, // AIService
			)

			err := src.Archive(
				context.Background(),
				testCase.managerID,
				testCase.planID,
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
