package plan

import (
	"context"
	"core_service/internal/domain"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlanService_ListByManager(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
	)

	managerID := uuid.New()
	employeeID1 := uuid.New()
	employeeID2 := uuid.New()
	plan1ID := uuid.New()
	plan2ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		managerID uuid.UUID

		mockBehavior mockBehavior

		expectedPlans    []*domain.Plan
		expectedErr      error
		expectedContains string
	}{
		{
			name:      "Успешное получение списка планов менеджера",
			managerID: managerID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				plans := []*domain.Plan{
					{
						ID:               plan1ID,
						EmployeeID:       employeeID1,
						CreatedBy:        managerID,
						Title:            "Plan 1",
						Description:      strPtr("Plan 1 Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					{
						ID:               plan2ID,
						EmployeeID:       employeeID2,
						CreatedBy:        managerID,
						Title:            "Plan 2",
						Description:      strPtr("Plan 2 Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         80,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
				}

				planRepo.EXPECT().
					ListByManager(gomock.Any(), managerID).
					Return(plans, nil)
			},

			expectedPlans: []*domain.Plan{
				{
					ID:               plan1ID,
					EmployeeID:       employeeID1,
					CreatedBy:        managerID,
					Title:            "Plan 1",
					Description:      strPtr("Plan 1 Description"),
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         50,
					Status:           domain.PlanActive,
					CreatedAt:        now,
					UpdatedAt:        now,
				},
				{
					ID:               plan2ID,
					EmployeeID:       employeeID2,
					CreatedBy:        managerID,
					Title:            "Plan 2",
					Description:      strPtr("Plan 2 Description"),
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         80,
					Status:           domain.PlanActive,
					CreatedAt:        now,
					UpdatedAt:        now,
				},
			},
			expectedErr: nil,
		},
		{
			name:      "Пустой список планов",
			managerID: managerID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ListByManager(gomock.Any(), managerID).
					Return([]*domain.Plan{}, nil)
			},

			expectedPlans: []*domain.Plan{},
			expectedErr:   nil,
		},
		{
			name:      "Ошибка при получении планов",
			managerID: managerID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ListByManager(gomock.Any(), managerID).
					Return(nil, errors.New("database error"))
			},

			expectedPlans:    nil,
			expectedContains: "list plans by manager: database error",
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

			result, err := src.ListByManager(
				context.Background(),
				testCase.managerID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, result)
			} else if testCase.expectedContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedContains)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testCase.expectedPlans), len(result))
				if len(result) > 0 {
					assert.Equal(t, testCase.expectedPlans[0].ID, result[0].ID)
					assert.Equal(t, testCase.expectedPlans[0].Title, result[0].Title)
					assert.Equal(t, testCase.expectedPlans[0].CreatedBy, result[0].CreatedBy)
				}
			}
		})
	}
}
