package auth

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/repository/redis"
	mock_redis "core_service/internal/repository/redis/mocks"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuthService_Refresh(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		authRepo *mock_postgres.MockAuthRepository,
		sessionRepo *mock_redis.MockSessionRepository,
	)

	userID := uuid.New()
	jwtService := jwt.NewJWTService(testJWTConfig())

	refreshToken, refreshJTI, err := jwtService.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	refreshClaims, err := jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		ID:    userID,
		Email: "test@mail.ru",
		Role:  domain.RoleEmployee,
	}

	testTable := []struct {
		name string

		refreshToken string

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешное обновление (сессия существует в Redis)",

			refreshToken: refreshToken,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				session := &redis.Session{
					UserID:    userID,
					TokenHash: hashRefreshToken(refreshToken),
				}

				sessionRepo.EXPECT().
					GetSession(gomock.Any(), refreshJTI).
					Return(session, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(nil)

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
			name: "Успешное обновление (сессия не в Redis, берем из БД)",

			refreshToken: refreshToken,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				sessionRepo.EXPECT().
					GetSession(gomock.Any(), refreshJTI).
					Return(nil, redis.ErrSessionNotFound)

				token := &domain.RefreshToken{
					UserID:    userID,
					JTI:       refreshJTI,
					TokenHash: hashRefreshToken(refreshToken),
					ExpiresAt: time.Now().Add(time.Hour),
					Revoked:   false,
				}

				authRepo.EXPECT().
					GetRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(token, nil)

				sessionRepo.EXPECT().
					SaveSession(gomock.Any(), refreshJTI, userID, hashRefreshToken(refreshToken), gomock.Any()).
					Return(nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				authRepo.EXPECT().
					DeleteRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil)

				sessionRepo.EXPECT().
					DeleteSession(gomock.Any(), refreshClaims.ID).
					Return(nil)

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
			name: "Неверный refresh токен",

			refreshToken: "invalid-token",

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
			},

			expectedErr: ErrInvalidRefreshToken,
		},
		{
			name: "Refresh токен не найден в БД и Redis",

			refreshToken: refreshToken,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				sessionRepo.EXPECT().
					GetSession(gomock.Any(), refreshJTI).
					Return(nil, redis.ErrSessionNotFound)

				authRepo.EXPECT().
					GetRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(nil, postgres.ErrRefreshTokenNotFound)
			},

			expectedErr: ErrInvalidRefreshToken,
		},
		{
			name: "Refresh токен revoked",

			refreshToken: refreshToken,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				sessionRepo.EXPECT().
					GetSession(gomock.Any(), refreshJTI).
					Return(nil, redis.ErrSessionNotFound)

				token := &domain.RefreshToken{
					UserID:    userID,
					JTI:       refreshJTI,
					TokenHash: hashRefreshToken(refreshToken),
					ExpiresAt: time.Now().Add(time.Hour),
					Revoked:   true,
				}

				authRepo.EXPECT().
					GetRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(token, nil)
			},

			expectedErr: ErrInvalidRefreshToken,
		},
		{
			name: "Refresh токен истек",

			refreshToken: refreshToken,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				sessionRepo.EXPECT().
					GetSession(gomock.Any(), refreshJTI).
					Return(nil, redis.ErrSessionNotFound)

				token := &domain.RefreshToken{
					UserID:    userID,
					JTI:       refreshJTI,
					TokenHash: hashRefreshToken(refreshToken),
					ExpiresAt: time.Now().Add(-time.Hour),
					Revoked:   false,
				}

				authRepo.EXPECT().
					GetRefreshToken(gomock.Any(), refreshClaims.ID).
					Return(token, nil)
			},

			expectedErr: ErrInvalidRefreshToken,
		},
		{
			name: "Хэш токена не совпадает",

			refreshToken: refreshToken,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				authRepo *mock_postgres.MockAuthRepository,
				sessionRepo *mock_redis.MockSessionRepository,
			) {
				session := &redis.Session{
					UserID:    userID,
					TokenHash: "wrong-hash",
				}

				sessionRepo.EXPECT().
					GetSession(gomock.Any(), refreshJTI).
					Return(session, nil)
			},

			expectedErr: ErrInvalidRefreshToken,
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

			src := NewAuthService(authRepo, userRepo, jwtService, sessionRepo)

			newAccess, newRefresh, err := src.Refresh(
				context.Background(),
				testCase.refreshToken,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Empty(t, newAccess)
				assert.Empty(t, newRefresh)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newAccess)
				assert.NotEmpty(t, newRefresh)
			}
		})
	}
}
