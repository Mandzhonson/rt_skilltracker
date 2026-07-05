package router

import (
	"core_service/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.CreateUser)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			protected := auth.Group("")
			protected.Use(authMiddleware)
			{
				protected.POST("/logout", authHandler.Logout)
			}
		}
		user := api.Group("/users")
		{
			user.Use(authMiddleware)
			user.GET("/me", userHandler.GetProfile)
			user.PATCH("/me", userHandler.UpdateProfile)
			user.PATCH("/me/password") // TODO
		}
	}

	return router
}
