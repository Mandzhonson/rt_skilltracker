package http

import (
	"core_service/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handler.AuthHandler,
) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.CreateUser)
		}
	}

	return router
}
