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

func TestUserService_CreateAdmin(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	email := "admin@mail.ru"
	password := "password123"
	firstName := "Admin"
	lastName := "User"
	userID := uuid.New()

	testTable := []struct {
		name string

		input CreateUserInput

		mockBehavior mockBehavior

		expectedErr      error
		expectedContains string
	}{
		{
			name: "Успешное создание администратора (админов нет)",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Role:      domain.RoleAdmin,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					ExistsAdmin(gomock.Any()).
					Return(false, nil)

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, postgres.ErrUserNotFound)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(userID, nil)
			},

			expectedErr: nil,
		},
		{
			name: "Администратор уже существует - пропускаем создание",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Role:      domain.RoleAdmin,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					ExistsAdmin(gomock.Any()).
					Return(true, nil)
			},

			expectedErr: nil,
		},
		{
			name: "Ошибка при проверке существования администратора",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Role:      domain.RoleAdmin,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					ExistsAdmin(gomock.Any()).
					Return(false, errors.New("database error"))
			},

			expectedErr: errors.New("database error"),
		},
		{
			name: "Ошибка при создании администратора",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Role:      domain.RoleAdmin,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					ExistsAdmin(gomock.Any()).
					Return(false, nil)

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, postgres.ErrUserNotFound)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("create error"))
			},

			expectedContains: "create user: create error",
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

			err := src.CreateAdmin(
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
