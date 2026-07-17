package user

import (
	"context"
	"core_service/internal/domain"
	mock_minio "core_service/internal/repository/minio/mocks"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_UpdatePassword(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	oldPassword := "oldPassword123"
	newPassword := "newPassword123"
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	testTable := []struct {
		name string

		userID      uuid.UUID
		oldPassword string
		newPassword string

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:        "Успешное обновление пароля",
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: newPassword,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:           userID,
					Email:        "test@mail.ru",
					PasswordHash: string(hashedOldPassword),
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					UpdatePassword(gomock.Any(), userID, gomock.Any()).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:        "Новый пароль слишком короткий",
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: "123",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedErr: ErrInvalidPassword,
		},
		{
			name:        "Пользователь не найден",
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: newPassword,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: ErrUserNotFound,
		},
		{
			name:        "Неверный старый пароль",
			userID:      userID,
			oldPassword: "wrongPassword",
			newPassword: newPassword,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:           userID,
					Email:        "test@mail.ru",
					PasswordHash: string(hashedOldPassword),
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)
			},

			expectedErr: ErrInvalidCredentials,
		},
		{
			name:        "Ошибка при обновлении пароля",
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: newPassword,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:           userID,
					Email:        "test@mail.ru",
					PasswordHash: string(hashedOldPassword),
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					UpdatePassword(gomock.Any(), userID, gomock.Any()).
					Return(errors.New("update error"))
			},

			expectedErr: errors.New("update error"),
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

			err := src.UpdatePassword(
				context.Background(),
				testCase.userID,
				testCase.oldPassword,
				testCase.newPassword,
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
