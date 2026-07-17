package user

import (
	"context"
	"core_service/internal/domain"
	mock_minio "core_service/internal/repository/minio/mocks"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetEmployeeAvatar(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		storage *mock_minio.MockStorage,
	)

	employeeID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	avatarKey := "avatars/employee_123.jpg"
	contentType := "image/jpeg"

	testTable := []struct {
		name string

		employeeID uuid.UUID
		managerID  uuid.UUID

		mockBehavior mockBehavior

		expectedReader  io.ReadCloser
		expectedContent string
		expectedErr     error
	}{
		{
			name:       "Успешное получение аватара сотрудника",
			employeeID: employeeID,
			managerID:  managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					AvatarKey: &avatarKey,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

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
			name:       "Сотрудник не найден",
			employeeID: employeeID,
			managerID:  managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     ErrUserNotFound,
		},
		{
			name:       "У сотрудника нет менеджера",
			employeeID: employeeID,
			managerID:  managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					AvatarKey: &avatarKey,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)
			},

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     ErrEmployeeNotAssigned,
		},
		{
			name:       "Сотрудник прикреплен к другому менеджеру",
			employeeID: employeeID,
			managerID:  otherManagerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					AvatarKey: &avatarKey,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)
			},

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     ErrEmployeeNotAssigned,
		},
		{
			name:       "У сотрудника нет аватара",
			employeeID: employeeID,
			managerID:  managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					AvatarKey: nil,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)
			},

			expectedReader:  nil,
			expectedContent: "",
			expectedErr:     ErrAvatarNotFound,
		},
		{
			name:       "Ошибка при получении аватара из хранилища",
			employeeID: employeeID,
			managerID:  managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				storage *mock_minio.MockStorage,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					AvatarKey: &avatarKey,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

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

			reader, contentType, err := src.GetEmployeeAvatar(
				context.Background(),
				testCase.employeeID,
				testCase.managerID,
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
