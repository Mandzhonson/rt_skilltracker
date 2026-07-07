package router

import (
	"core_service/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc) *gin.Engine {
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
			user.PATCH("/me/password", userHandler.UpdatePassword)

			user.PUT("/me/avatar", userHandler.SetAvatar)
			user.GET("/me/avatar", userHandler.GetAvatar)
			user.DELETE("/me/avatar", userHandler.DeleteAvatar)
		}
		admin := api.Group("/admin")
		{
			admin.Use(authMiddleware, adminMiddleware)
			admin.GET("/users", adminHandler.ListUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.PATCH("/users/:id/role", adminHandler.UpdateRole)
			admin.PATCH("/users/:id/manager", adminHandler.AssignManager)
			admin.DELETE("/users/:id/manager", adminHandler.RemoveManager)
			admin.GET("/managers/:id/employees", adminHandler.ListEmployeesByManager)

		}
	}

	return router
}
