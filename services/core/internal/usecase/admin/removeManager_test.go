package admin

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

func TestAdminService_RemoveManager(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		userID uuid.UUID

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:   "Успешное удаление менеджера",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					RemoveManager(gomock.Any(), userID).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:   "Пользователь не найден",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
		},
		{
			name:   "Менеджер не назначен",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)
			},

			expectedErr: ErrManagerNotAssigned,
		},
		{
			name:   "Ошибка при удалении менеджера",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					RemoveManager(gomock.Any(), userID).
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

			err := src.RemoveManager(
				context.Background(),
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