package plan

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

func TestPlanService_Update(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
	)

	planID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	title := "Updated Plan Title"
	description := "Updated Description"

	testTable := []struct {
		name string

		input UpdatePlanInput

		mockBehavior mockBehavior

		expectedErr      error
		expectedContains string
	}{
		{
			name: "Успешное обновление плана",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Old Title",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				planRepo.EXPECT().
					Update(gomock.Any(), planID, title, &description).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное обновление плана (без описания)",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       title,
				Description: nil,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Old Title",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				planRepo.EXPECT().
					Update(gomock.Any(), planID, title, nil).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Пустое название (только пробелы)",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       "   ",
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidTitle,
		},
		{
			name: "Название меньше 3 символов",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       "Ab",
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedErr: ErrInvalidTitle,
		},
		{
			name: "План не найден",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       title,
				Description: &description,
			},

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
			name: "Нет прав (другой менеджер)",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   otherManagerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Old Title",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedErr: ErrManagerForbidden,
		},
		{
			name: "Ошибка при получении плана",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, errors.New("database error"))
			},

			expectedContains: "get plan: database error",
		},
		{
			name: "Ошибка при обновлении плана",
			input: UpdatePlanInput{
				PlanID:      planID,
				ManagerID:   managerID,
				Title:       title,
				Description: &description,
			},

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
			) {
				plan := &domain.Plan{
					ID:        planID,
					CreatedBy: managerID,
					Title:     "Old Title",
					Status:    domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				planRepo.EXPECT().
					Update(gomock.Any(), planID, title, &description).
					Return(errors.New("update error"))
			},

			expectedContains: "update plan: update error",
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

			err := src.Update(
				context.Background(),
				testCase.input,
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
