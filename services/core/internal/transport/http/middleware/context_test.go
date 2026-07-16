package middleware

import (
	"core_service/internal/pkg/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testClaims := &jwt.Claims{
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject: uuid.New().String(),
		},
	}

	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("claims", testClaims)
		c.Next()
	})

	r.GET("/test", func(c *gin.Context) {
		claims, ok := GetClaims(c)
		assert.True(t, ok)
		assert.NotNil(t, claims)
		assert.Equal(t, testClaims, claims)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUserID := uuid.New()

	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("userID", testUserID)
		c.Next()
	})

	r.GET("/test", func(c *gin.Context) {
		id, ok := GetUserID(c)
		assert.True(t, ok)
		assert.Equal(t, testUserID, id)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
