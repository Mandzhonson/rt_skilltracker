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

func TestUserService_CreateUser(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	email := "test@mail.ru"
	password := "password123"
	firstName := "Test"
	lastName := "User"
	position := "Developer"
	userID := uuid.New()

	testTable := []struct {
		name string

		input CreateUserInput

		mockBehavior mockBehavior

		expectedID       uuid.UUID
		expectedErr      error
		expectedContains string
	}{
		{
			name: "Успешное создание пользователя",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, postgres.ErrUserNotFound)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(userID, nil)
			},

			expectedID:  userID,
			expectedErr: nil,
		},
		{
			name: "Успешное создание пользователя (роль не указана - по умолчанию employee)",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      "",
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, postgres.ErrUserNotFound)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(userID, nil)
			},

			expectedID:  userID,
			expectedErr: nil,
		},
		{
			name: "Неверный email",
			input: CreateUserInput{
				Email:     "invalid-email",
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidEmail,
		},
		{
			name: "Слишком короткий пароль",
			input: CreateUserInput{
				Email:     email,
				Password:  "123",
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidPassword,
		},
		{
			name: "Пустое имя",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: "   ",
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidName,
		},
		{
			name: "Пустая фамилия",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  "   ",
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrInvalidName,
		},
		{
			name: "Пользователь уже существует",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				existingUser := &domain.User{
					ID:    uuid.New(),
					Email: email,
				}

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(existingUser, nil)
			},

			expectedID:  uuid.Nil,
			expectedErr: ErrUserAlreadyExists,
		},
		{
			name: "Ошибка при проверке существования пользователя",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, errors.New("database error"))
			},

			expectedID:       uuid.Nil,
			expectedContains: "check user exists: database error",
		},
		{
			name: "Ошибка при создании пользователя",
			input: CreateUserInput{
				Email:     email,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
				Position:  position,
				Role:      domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, postgres.ErrUserNotFound)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("create error"))
			},

			expectedID:       uuid.Nil,
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

			id, err := src.CreateUser(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Equal(t, testCase.expectedID, id)
			} else if testCase.expectedContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedContains)
				assert.Equal(t, testCase.expectedID, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedID, id)
			}
		})
	}
}
