package middleware

import (
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeMiddleware(t *testing.T) {
	testTable := []struct {
		name string

		role domain.Role

		expectedStatus int
	}{
		{
			name:           "Сотрудник - доступ разрешен",
			role:           domain.RoleEmployee,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Менеджер - доступ запрещен",
			role:           domain.RoleManager,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Администратор - доступ запрещен",
			role:           domain.RoleAdmin,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Нет claims - доступ запрещен",
			role:           "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				if testCase.role != "" {
					claims := &jwt.Claims{
						RegisteredClaims: jwtv5.RegisteredClaims{
							Subject: uuid.New().String(),
						},
						Role: testCase.role,
					}
					c.Set("claims", claims)
				}
				c.Next()
			})
			r.Use(EmployeeMiddleware())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
