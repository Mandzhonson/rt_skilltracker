package middleware

import (
	"core_service/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetClaims(c *gin.Context) (*jwt.Claims, bool) {
	claims, ok := c.Get("claims")
	if !ok {
		return nil, false
	}

	result, ok := claims.(*jwt.Claims)
	if !ok {
		return nil, false
	}

	return result, true
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return userID, true
}
