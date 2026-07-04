package router

import (
	"core_service/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.CreateUser)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			protected := auth.Group("")
			protected.Use(authMiddleware)
			{
				protected.POST("/logout", authHandler.Logout)
			}
		}
	}

	return router
}
