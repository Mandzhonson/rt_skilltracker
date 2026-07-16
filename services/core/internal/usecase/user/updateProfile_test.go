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

func TestUserService_UpdateProfile(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	email := "updated@mail.ru"
	firstName := "Updated"
	lastName := "User"
	emptyString := ""

	testTable := []struct {
		name string

		upd UpdateProfileInput

		mockBehavior mockBehavior

		expectedUser *domain.User
		expectedErr  error
	}{
		{
			name: "Успешное обновление всех полей",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: &firstName,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdateProfile(gomock.Any(), userID, gomock.Any()).
					Return(nil)

				updatedUser := &domain.User{
					ID:        userID,
					Email:     email,
					FirstName: firstName,
					LastName:  lastName,
					Role:      domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(updatedUser, nil)
			},

			expectedUser: &domain.User{
				ID:        userID,
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
				Role:      domain.RoleEmployee,
			},
			expectedErr: nil,
		},
		{
			name: "Успешное обновление только email",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: nil,
				LastName:  nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdateProfile(gomock.Any(), userID, gomock.Any()).
					Return(nil)

				updatedUser := &domain.User{
					ID:        userID,
					Email:     email,
					FirstName: "Old",
					LastName:  "User",
					Role:      domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(updatedUser, nil)
			},

			expectedUser: &domain.User{
				ID:        userID,
				Email:     email,
				FirstName: "Old",
				LastName:  "User",
				Role:      domain.RoleEmployee,
			},
			expectedErr: nil,
		},
		{
			name: "Успешное обновление только имени",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     nil,
				FirstName: &firstName,
				LastName:  nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdateProfile(gomock.Any(), userID, gomock.Any()).
					Return(nil)

				updatedUser := &domain.User{
					ID:        userID,
					Email:     "old@mail.ru",
					FirstName: firstName,
					LastName:  "User",
					Role:      domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(updatedUser, nil)
			},

			expectedUser: &domain.User{
				ID:        userID,
				Email:     "old@mail.ru",
				FirstName: firstName,
				LastName:  "User",
				Role:      domain.RoleEmployee,
			},
			expectedErr: nil,
		},
		{
			name: "Нет данных для обновления",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     nil,
				FirstName: nil,
				LastName:  nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedUser: nil,
			expectedErr:  ErrNoContent,
		},
		{
			name: "Неверный email (пустая строка)",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &emptyString,
				FirstName: &firstName,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedUser: nil,
			expectedErr:  ErrInvalidEmail,
		},
		{
			name: "Неверный email (невалидный формат)",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     strPtr("invalid-email"),
				FirstName: &firstName,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedUser: nil,
			expectedErr:  ErrInvalidEmail,
		},
		{
			name: "Пустое имя",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: &emptyString,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedUser: nil,
			expectedErr:  ErrInvalidName,
		},
		{
			name: "Пустая фамилия",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: &firstName,
				LastName:  &emptyString,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedUser: nil,
			expectedErr:  ErrInvalidName,
		},
		{
			name: "Пользователь уже существует (email занят)",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: &firstName,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdateProfile(gomock.Any(), userID, gomock.Any()).
					Return(postgres.ErrUserAlreadyExists)
			},

			expectedUser: nil,
			expectedErr:  ErrUserAlreadyExists,
		},
		{
			name: "Пользователь не найден",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: &firstName,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdateProfile(gomock.Any(), userID, gomock.Any()).
					Return(postgres.ErrUserNotFound)
			},

			expectedUser: nil,
			expectedErr:  ErrUserNotFound,
		},
		{
			name: "Ошибка при обновлении профиля",
			upd: UpdateProfileInput{
				UserID:    userID,
				Email:     &email,
				FirstName: &firstName,
				LastName:  &lastName,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					UpdateProfile(gomock.Any(), userID, gomock.Any()).
					Return(errors.New("database error"))
			},

			expectedUser: nil,
			expectedErr:  errors.New("database error"),
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

			user, err := src.UpdateProfile(
				context.Background(),
				testCase.upd,
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
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
