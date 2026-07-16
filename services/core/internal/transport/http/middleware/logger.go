package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		latency := time.Since(start)

		status := c.Writer.Status()

		if status >= 500 {
			log.Error("request failed", slog.Int("status", status), slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path), slog.Duration("latency", latency), slog.String("ip", c.ClientIP()))
			return
		}

		if latency > time.Second {
			log.Warn("slow request", slog.Int("status", status), slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path), slog.Duration("latency", latency))
		}
	}
}

func Recovery(log *slog.Logger) gin.HandlerFunc {

	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered", slog.Any("panic", recovered), slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	})
}
