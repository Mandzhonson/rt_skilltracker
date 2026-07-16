package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	mock_redis "core_service/internal/repository/redis/mocks"
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		authRepo *mock_postgres.MockAuthRepository,
		sessionRepo *mock_redis.MockSessionRepository,
	)

	testTable := []struct {
		name string

		email    string
		password string

		mockBehavior mockBehavior

		expectedErr      error
		expectedErrMatch string
	}{
		{
			name:     "Успешный вход",
			email:    "test@mail.ru",
			password: "password123",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

				user := &domain.User{
					ID:           uuid.New(),
					Email:        "test@mail.ru",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
					Role:         domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@mail.ru").
					Return(user, nil)

				authRepo.EXPECT().
					SaveRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil)

				sessionRepo.EXPECT().
					SaveSession(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:     "Пользователь не найден",
			email:    "notfound@mail.ru",
			password: "password123",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "notfound@mail.ru").
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: ErrInvalidCredentials,
		},
		{
			name:     "Неверный пароль",
			email:    "test@mail.ru",
			password: "wrongpassword",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

				user := &domain.User{
					ID:           uuid.New(),
					Email:        "test@mail.ru",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
					Role:         domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@mail.ru").
					Return(user, nil)
			},

			expectedErr: ErrInvalidCredentials,
		},
		{
			name:     "Ошибка при сохранении refresh токена",
			email:    "test@mail.ru",
			password: "password123",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

				user := &domain.User{
					ID:           uuid.New(),
					Email:        "test@mail.ru",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
					Role:         domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@mail.ru").
					Return(user, nil)

				authRepo.EXPECT().
					SaveRefreshToken(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},

			expectedErrMatch: "save refresh token: database error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			authRepo := mock_postgres.NewMockAuthRepository(ctrl)
			sessionRepo := mock_redis.NewMockSessionRepository(ctrl)

			testCase.mockBehavior(userRepo, authRepo, sessionRepo)

			jwtService := jwt.NewJWTService(testJWTConfig())
			src := NewAuthService(authRepo, userRepo, jwtService, sessionRepo)

			accessToken, refreshToken, err := src.Login(
				context.Background(),
				testCase.email,
				testCase.password,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
			} else if testCase.expectedErrMatch != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedErrMatch)
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
			}
		})
	}
}

func TestAuthService_authenticateUser(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	testTable := []struct {
		name string

		email    string
		password string

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:     "Успешная аутентификация",
			email:    "test@mail.ru",
			password: "password123",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

				user := &domain.User{
					ID:           uuid.New(),
					Email:        "test@mail.ru",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
					Role:         domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@mail.ru").
					Return(user, nil)
			},

			expectedErr: nil,
		},
		{
			name:     "Пользователь не найден",
			email:    "notfound@mail.ru",
			password: "password123",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "notfound@mail.ru").
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: ErrInvalidCredentials,
		},
		{
			name:     "Неверный пароль",
			email:    "test@mail.ru",
			password: "wrongpassword",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

				user := &domain.User{
					ID:           uuid.New(),
					Email:        "test@mail.ru",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
					Role:         domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@mail.ru").
					Return(user, nil)
			},

			expectedErr: ErrInvalidCredentials,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			authRepo := mock_postgres.NewMockAuthRepository(ctrl)
			sessionRepo := mock_redis.NewMockSessionRepository(ctrl)

			testCase.mockBehavior(userRepo)

			jwtService := jwt.NewJWTService(testJWTConfig())
			src := NewAuthService(authRepo, userRepo, jwtService, sessionRepo)

			user, err := src.authenticateUser(
				context.Background(),
				testCase.email,
				testCase.password,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestAuthService_storeRefreshSession(t *testing.T) {
	type mockBehavior func(
		authRepo *mock_postgres.MockAuthRepository,
	)

	testTable := []struct {
		name string

		userID    uuid.UUID
		hashToken string
		jti       string

		mockBehavior mockBehavior

		expectedErr      error
		expectedErrMatch string
	}{
		{
			name:      "Успешное сохранение",
			userID:    uuid.New(),
			hashToken: "hashed-token",
			jti:       "jti-123",

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
			) {
				authRepo.EXPECT().
					SaveRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:      "Ошибка при сохранении",
			userID:    uuid.New(),
			hashToken: "hashed-token",
			jti:       "jti-123",

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
			) {
				authRepo.EXPECT().
					SaveRefreshToken(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},

			expectedErrMatch: "save refresh token: database error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			authRepo := mock_postgres.NewMockAuthRepository(ctrl)
			sessionRepo := mock_redis.NewMockSessionRepository(ctrl)

			testCase.mockBehavior(authRepo)

			jwtService := jwt.NewJWTService(testJWTConfig())
			src := NewAuthService(authRepo, userRepo, jwtService, sessionRepo)

			err := src.storeRefreshSession(
				context.Background(),
				testCase.userID,
				testCase.hashToken,
				testCase.jti,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else if testCase.expectedErrMatch != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedErrMatch)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
