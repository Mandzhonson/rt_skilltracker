package middleware

import (
	"core_service/internal/config"
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/redis/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	type mockBehavior func(s *mocks.MockSessionRepository)

	userID := uuid.New()
	jwtService := jwt.NewJWTService(config.JWTConfig{
		AccessSecret:  "test-secret",
		RefreshSecret: "test-secret",
		AccessTTL:     time.Hour,
		RefreshTTL:    time.Hour * 24,
	})

	validToken, err := jwtService.GenerateAccessToken(userID, "employee")
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name string

		authHeader string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:       "Успешная аутентификация",
			authHeader: "Bearer " + validToken,

			mockBehavior: func(s *mocks.MockSessionRepository) {
				s.EXPECT().
					IsBlacklisted(gomock.Any(), gomock.Any()).
					Return(false, nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:       "Отсутствует Authorization header",
			authHeader: "",

			mockBehavior: func(s *mocks.MockSessionRepository) {
			},

			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Неверный формат токена",
			authHeader: "Token invalid",

			mockBehavior: func(s *mocks.MockSessionRepository) {
			},

			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Неверный токен",
			authHeader: "Bearer invalid-token",

			mockBehavior: func(s *mocks.MockSessionRepository) {
			},

			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Токен в blacklist",
			authHeader: "Bearer " + validToken,

			mockBehavior: func(s *mocks.MockSessionRepository) {
				s.EXPECT().
					IsBlacklisted(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},

			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Ошибка при проверке blacklist",
			authHeader: "Bearer " + validToken,

			mockBehavior: func(s *mocks.MockSessionRepository) {
				s.EXPECT().
					IsBlacklisted(gomock.Any(), gomock.Any()).
					Return(false, errors.New("redis error"))
			},

			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionRepo := mocks.NewMockSessionRepository(ctrl)

			testCase.mockBehavior(sessionRepo)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(AuthMiddleware(jwtService, sessionRepo))
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if testCase.authHeader != "" {
				req.Header.Set("Authorization", testCase.authHeader)
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
