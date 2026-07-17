package handler

import (
	"bytes"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/usecase/auth"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthService)

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешный вход",
			body: `{"email":"test@mail.ru","password":"password123"}`,

			mockBehavior: func(s *mocks.MockAuthService) {
				s.EXPECT().
					Login(gomock.Any(), "test@mail.ru", "password123").
					Return("access-token", "refresh-token", nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"access-token","refresh_token":"refresh-token"}`,
		},
		{
			name: "Неверные учетные данные",
			body: `{"email":"test@mail.ru","password":"wrong"}`,

			mockBehavior: func(s *mocks.MockAuthService) {
				s.EXPECT().
					Login(gomock.Any(), "test@mail.ru", "wrong").
					Return("", "", auth.ErrInvalidCredentials)
			},

			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid email or password"}`,
		},
		{
			name: "Неверный запрос",
			body: `{"email":"test@mail.ru"`,

			mockBehavior: func(s *mocks.MockAuthService) {
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
		},
		{
			name: "Ошибка сервера",
			body: `{"email":"test@mail.ru","password":"password123"}`,

			mockBehavior: func(s *mocks.MockAuthService) {
				s.EXPECT().
					Login(gomock.Any(), "test@mail.ru", "password123").
					Return("", "", errors.New("database error"))
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAuthService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAuthHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/auth/login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthService)

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное обновление токенов",
			body: `{"refresh_token":"valid-refresh-token"}`,

			mockBehavior: func(s *mocks.MockAuthService) {
				s.EXPECT().
					Refresh(gomock.Any(), "valid-refresh-token").
					Return("new-access-token", "new-refresh-token", nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`,
		},
		{
			name: "Неверный refresh токен",
			body: `{"refresh_token":"invalid-token"}`,

			mockBehavior: func(s *mocks.MockAuthService) {
				s.EXPECT().
					Refresh(gomock.Any(), "invalid-token").
					Return("", "", auth.ErrInvalidRefreshToken)
			},

			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid refresh token"}`,
		},
		{
			name: "Неверный запрос",
			body: `{"refresh_token":"valid-token"`,

			mockBehavior: func(s *mocks.MockAuthService) {
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAuthService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAuthHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/auth/refresh", handler.Refresh)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthService)

	userID := uuid.New()
	jti := "test-jti"
	claims := &jwt.Claims{
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject: userID.String(),
			ID:      jti,
		},
	}

	testTable := []struct {
		name string

		body         string
		setupContext func(c *gin.Context)

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешный выход",
			body: `{"refresh_token":"valid-refresh-token"}`,

			setupContext: func(c *gin.Context) {
				c.Set("claims", claims)
			},

			mockBehavior: func(s *mocks.MockAuthService) {
				s.EXPECT().
					Logout(gomock.Any(), claims, "valid-refresh-token").
					Return(nil)
			},

			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name: "Неверный запрос",
			body: `{"refresh_token":"valid-token"`,

			setupContext: func(c *gin.Context) {
				c.Set("claims", claims)
			},

			mockBehavior: func(s *mocks.MockAuthService) {
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
		{
			name: "Неавторизован",
			body: `{"refresh_token":"valid-refresh-token"}`,

			setupContext: func(c *gin.Context) {
			},

			mockBehavior: func(s *mocks.MockAuthService) {
			},

			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAuthService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAuthHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				if testCase.setupContext != nil {
					testCase.setupContext(c)
				}
				c.Next()
			})
			r.POST("/auth/logout", handler.Logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.expectedStatus == http.StatusNoContent {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
		})
	}
}
