package admin

import (
	"context"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_UpdatePosition(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()

	testTable := []struct {
		name string

		userID   uuid.UUID
		position string

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:     "Успешное обновление должности",
			userID:   userID,
			position: "Senior Developer",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdatePosition(gomock.Any(), userID, "Senior Developer").
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:     "Успешное обновление должности (с пробелами в начале и конце)",
			userID:   userID,
			position: "  Team Lead  ",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdatePosition(gomock.Any(), userID, "  Team Lead  ").
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:     "Пустая должность (только пробелы) - ошибка",
			userID:   userID,
			position: "   ",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedErr: ErrInvalidPosition,
		},
		{
			name:     "Пустая должность (пустая строка) - ошибка",
			userID:   userID,
			position: "",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedErr: ErrInvalidPosition,
		},
		{
			name:     "Пользователь не найден",
			userID:   userID,
			position: "Developer",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdatePosition(gomock.Any(), userID, "Developer").
					Return(postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
		},
		{
			name:     "Ошибка при обновлении должности",
			userID:   userID,
			position: "Developer",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdatePosition(gomock.Any(), userID, "Developer").
					Return(assert.AnError)
			},

			expectedErr: assert.AnError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)

			testCase.mockBehavior(userRepo)

			src := NewAdminService(
				userRepo,
				nil, // planRepo
				nil, // skillRepo
				nil, // storage
			)

			err := src.UpdatePosition(
				context.Background(),
				testCase.userID,
				testCase.position,
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
