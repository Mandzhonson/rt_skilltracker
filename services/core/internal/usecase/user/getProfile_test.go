package user

import (
	"context"
	"core_service/internal/domain"
	mock_minio "core_service/internal/repository/minio/mocks"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetProfile(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		userID uuid.UUID

		mockBehavior mockBehavior

		expectedUser *domain.User
		expectedErr  error
	}{
		{
			name:   "Успешное получение профиля",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					FirstName: "Test",
					LastName:  "User",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)
			},

			expectedUser: &domain.User{
				ID:        userID,
				Email:     "test@mail.ru",
				FirstName: "Test",
				LastName:  "User",
				Role:      domain.RoleEmployee,
				Position:  "Developer",
				ManagerID: &managerID,
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

			expectedUser: nil,
			expectedErr:  ErrUserNotFound,
		},
		{
			name:   "Ошибка при получении пользователя",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(nil, assert.AnError)
			},

			expectedUser: nil,
			expectedErr:  assert.AnError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			storage := mock_minio.NewMockStorage(ctrl)

			testCase.mockBehavior(userRepo)

			src := NewUserService(
				userRepo,
				storage,
				skillRepo,
				planRepo,
			)

			user, err := src.GetProfile(
				context.Background(),
				testCase.userID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testCase.expectedUser.ID, user.ID)
				assert.Equal(t, testCase.expectedUser.Email, user.Email)
				assert.Equal(t, testCase.expectedUser.FirstName, user.FirstName)
				assert.Equal(t, testCase.expectedUser.LastName, user.LastName)
				assert.Equal(t, testCase.expectedUser.Role, user.Role)
				assert.Equal(t, testCase.expectedUser.Position, user.Position)
			}
		})
	}
}
