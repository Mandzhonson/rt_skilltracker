package middleware

import (
	"core_service/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*jwt.Claims, bool) {
	value, ok := c.Get("claims")
	if !ok {
		return nil, false
	}

	claims, ok := value.(*jwt.Claims)
	if !ok {
		return nil, false
	}

	return claims, true
}
