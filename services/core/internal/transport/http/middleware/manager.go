package middleware

import (
	"core_service/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ManagerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if claims.Role != domain.RoleManager {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
