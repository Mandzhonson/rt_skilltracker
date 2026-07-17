package middleware

import (
	"core_service/internal/pkg/jwt"
	"core_service/internal/repository/redis"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(jc *jwt.JWTService, blackListService redis.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		const bearerPrefix string = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := jc.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ok, err := blackListService.IsBlacklisted(c.Request.Context(), claims.ID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("claims", claims)
		c.Set("userID", userID)

		c.Next()
	}
}
