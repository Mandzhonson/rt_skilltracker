package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/repository/redis"
	mock_redis "core_service/internal/repository/redis/mocks"
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuthService_Logout(t *testing.T) {
	type mockBehavior func(
		authRepo *mock_postgres.MockAuthRepository,
		sessionRepo *mock_redis.MockSessionRepository,
	)

	userID := uuid.New()
	jwtService := jwt.NewJWTService(testJWTConfig())

	accessToken, err := jwtService.GenerateAccessToken(userID, domain.RoleEmployee)
	if err != nil {
		t.Fatal(err)
	}

	accessClaims, err := jwtService.ParseAccessToken(accessToken)
	if err != nil {
		t.Fatal(err)
	}

	refreshToken, jti, err := jwtService.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	refreshClaims, err := jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		t.Fatal(err)
	}
	_ = jti

	testTable := []struct {
		name string

		claimsAccess *jwt.Claims
		refreshToken string

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешный выход",

			claimsAccess: accessClaims,
			refreshToken: refreshToken,

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					BlacklistAccessToken(gomock.Any(), accessClaims.ID, gomock.Any()).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Неверный refresh токен",

			claimsAccess: accessClaims,
			refreshToken: "invalid-token",

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
			},

			expectedErr: jwt.ErrTokenInvalid,
		},
		{
			name: "Subject не совпадает",

			claimsAccess: accessClaims,
			refreshToken: func() string {
				otherUserID := uuid.New()
				token, _, err := jwtService.GenerateRefreshToken(otherUserID)
				if err != nil {
					t.Fatal(err)
				}
				return token
			}(),

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
			},

			expectedErr: ErrInvalidCredentials,
		},
		{
			name: "Ошибка при удалении refresh токена (игнорируем ErrRefreshTokenNotFound)",

			claimsAccess: accessClaims,
			refreshToken: refreshToken,

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(postgres.ErrRefreshTokenNotFound)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					BlacklistAccessToken(gomock.Any(), accessClaims.ID, gomock.Any()).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Ошибка при удалении refresh токена (не игнорируем)",

			claimsAccess: accessClaims,
			refreshToken: refreshToken,

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(errors.New("database error"))
			},

			expectedErr: errors.New("database error"),
		},
		{
			name: "Ошибка при удалении сессии (игнорируем ErrSessionNotFound)",

			claimsAccess: accessClaims,
			refreshToken: refreshToken,

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(redis.ErrSessionNotFound)

				sessionRepo.EXPECT().
					BlacklistAccessToken(gomock.Any(), accessClaims.ID, gomock.Any()).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Ошибка при удалении сессии (не игнорируем)",

			claimsAccess: accessClaims,
			refreshToken: refreshToken,

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(errors.New("redis error"))
			},

			expectedErr: errors.New("redis error"),
		},
		{
			name: "Ошибка при добавлении в blacklist",

			claimsAccess: accessClaims,
			refreshToken: refreshToken,

			mockBehavior: func(
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					BlacklistAccessToken(gomock.Any(), accessClaims.ID, gomock.Any()).
					Return(errors.New("redis blacklist error"))
			},

			expectedErr: errors.New("redis blacklist error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			authRepo := mock_postgres.NewMockAuthRepository(ctrl)
			sessionRepo := mock_redis.NewMockSessionRepository(ctrl)

			testCase.mockBehavior(authRepo, sessionRepo)

			src := NewAuthService(authRepo, userRepo, jwtService, sessionRepo)

			err := src.Logout(
				context.Background(),
				testCase.claimsAccess,
				testCase.refreshToken,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
