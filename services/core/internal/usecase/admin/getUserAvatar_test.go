package admin

import (
	"context"
	"core_service/internal/domain"
	mock_minio "core_service/internal/repository/minio/mocks"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/user"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_GetUserAvatar(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		storage *mock_minio.MockStorage,
	)

	userID := uuid.New()
	avatarKey := "avatars/user_123.jpg"
	contentType := "image/jpeg"

	testTable := []struct {
		name string

		userID uuid.UUID

		mockBehavior mockBehavior

		expectedReader  io.ReadCloser
		expectedContent string
		expectedErr     error
	}{
		{
			name:   "Успешное получение аватара",
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

				mockReader := &mockReadCloser{}

				storage.EXPECT().
					GetAvatar(gomock.Any(), avatarKey).
					Return(mockReader, contentType, nil)
			},

			expectedReader:  &mockReadCloser{},
			expectedContent: contentType,
			expectedErr:     nil,
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

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     postgres.ErrUserNotFound,
		},
		{
			name:   "У пользователя нет аватара",
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

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     user.ErrAvatarNotFound,
		},
		{
			name:   "Ошибка при получении аватара из хранилища",
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
					GetAvatar(gomock.Any(), avatarKey).
					Return(nil, "", errors.New("storage error"))
			},

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     errors.New("storage error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			storage := mock_minio.NewMockStorage(ctrl)

			testCase.mockBehavior(userRepo, storage)

			src := NewAdminService(
				userRepo,
				nil, // planRepo
				nil, // skillRepo
				storage,
			)

			reader, contentType, err := src.GetUserAvatar(
				context.Background(),
				testCase.userID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, reader)
				assert.Empty(t, contentType)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				assert.Equal(t, testCase.expectedContent, contentType)
			}
		})
	}
}

type mockReadCloser struct{}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockReadCloser) Close() error {
	return nil
}
