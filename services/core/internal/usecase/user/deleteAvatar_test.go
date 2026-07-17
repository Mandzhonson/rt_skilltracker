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
)

func TestUserService_DeleteAvatar(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		storage *mock_minio.MockStorage,
	)

	userID := uuid.New()
	avatarKey := "avatars/user_123.jpg"

	testTable := []struct {
		name string

		userID uuid.UUID

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:   "Успешное удаление аватара",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					AvatarKey: &avatarKey,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				storage.EXPECT().
					DeleteAvatar(gomock.Any(), avatarKey).
					Return(nil)

				userRepo.EXPECT().
					UpdateAvatar(gomock.Any(), userID, nil).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:   "Пользователь не найден",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
		},
		{
			name:   "У пользователя нет аватара - ErrNoContent",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					AvatarKey: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)
			},

			expectedErr: ErrNoContent,
		},
		{
			name:   "Ошибка при удалении аватара из хранилища",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					AvatarKey: &avatarKey,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				storage.EXPECT().
					DeleteAvatar(gomock.Any(), avatarKey).
					Return(errors.New("storage error"))
			},

			expectedErr: errors.New("storage error"),
		},
		{
			name:   "Ошибка при обновлении пользователя (удаление avatar_key)",
			userID: userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					AvatarKey: &avatarKey,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				storage.EXPECT().
					DeleteAvatar(gomock.Any(), avatarKey).
					Return(nil)

				userRepo.EXPECT().
					UpdateAvatar(gomock.Any(), userID, nil).
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

			testCase.mockBehavior(userRepo, storage)

			src := NewUserService(
				userRepo,
				storage,
				skillRepo,
				planRepo,
			)

			err := src.DeleteAvatar(
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
