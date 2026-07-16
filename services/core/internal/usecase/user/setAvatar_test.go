package user

import (
	"context"
	"core_service/internal/domain"
	mock_minio "core_service/internal/repository/minio/mocks"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_SetAvatar(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		storage *mock_minio.MockStorage,
	)

	userID := uuid.New()
	objectKey := userID.String() + ".jpg"
	oldAvatarKey := "avatars/old_avatar.jpg"
	contentType := "image/jpeg"
	fileSize := int64(1024 * 1024)

	testTable := []struct {
		name string

		input SetAvatarInput

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешная загрузка аватара (без старого аватара)",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        fileSize,
				ContentType: contentType,
			},

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

				storage.EXPECT().
					UploadAvatar(gomock.Any(), objectKey, gomock.Any(), fileSize, contentType).
					Return(objectKey, nil)

				userRepo.EXPECT().
					UpdateAvatar(gomock.Any(), userID, &objectKey).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешная загрузка аватара (с удалением старого)",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        fileSize,
				ContentType: contentType,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					AvatarKey: &oldAvatarKey,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				storage.EXPECT().
					UploadAvatar(gomock.Any(), objectKey, gomock.Any(), fileSize, contentType).
					Return(objectKey, nil)

				userRepo.EXPECT().
					UpdateAvatar(gomock.Any(), userID, &objectKey).
					Return(nil)

				storage.EXPECT().
					DeleteAvatar(gomock.Any(), oldAvatarKey).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Файл слишком большой",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        maxAvatarSize + 1,
				ContentType: contentType,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
			},

			expectedErr: ErrAvatarTooLarge,
		},
		{
			name: "Неверный формат файла",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        fileSize,
				ContentType: "image/gif",
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
			},

			expectedErr: ErrInvalidAvatarFormat,
		},
		{
			name: "Пользователь не найден",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        fileSize,
				ContentType: contentType,
			},

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
			name: "Ошибка при загрузке аватара в хранилище",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        fileSize,
				ContentType: contentType,
			},

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

				storage.EXPECT().
					UploadAvatar(gomock.Any(), objectKey, gomock.Any(), fileSize, contentType).
					Return("", errors.New("upload error"))
			},

			expectedErr: errors.New("upload error"),
		},
		{
			name: "Ошибка при обновлении пользователя (удаляем загруженный файл)",
			input: SetAvatarInput{
				UserID:      userID,
				File:        io.NopCloser(strings.NewReader("test")),
				Size:        fileSize,
				ContentType: contentType,
			},

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

				storage.EXPECT().
					UploadAvatar(gomock.Any(), objectKey, gomock.Any(), fileSize, contentType).
					Return(objectKey, nil)

				userRepo.EXPECT().
					UpdateAvatar(gomock.Any(), userID, &objectKey).
					Return(errors.New("update error"))

				storage.EXPECT().
					DeleteAvatar(gomock.Any(), objectKey).
					Return(nil)
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

			err := src.SetAvatar(
				context.Background(),
				testCase.input,
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
