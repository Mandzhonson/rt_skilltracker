package router

import (
	"core_service/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	planHandler *handler.PlanHandler,
	taskHandler *handler.TaskHandler,
	testHandler *handler.TestHandler,
	skillHandler *handler.SkillHandler,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
	managerMiddleware gin.HandlerFunc,
	employeeMiddleware gin.HandlerFunc) *gin.Engine {
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
		authorized := api.Group("")
		authorized.Use(authMiddleware)
		user := authorized.Group("/users")
		{
			user.GET("/me", userHandler.GetProfile)
			user.PATCH("/me", userHandler.UpdateProfile)
			user.PATCH("/me/password", userHandler.UpdatePassword)
			user.PUT("/me/avatar", userHandler.SetAvatar)
			user.GET("/me/avatar", userHandler.GetAvatar)
			user.DELETE("/me/avatar", userHandler.DeleteAvatar)
		}
		admin := authorized.Group("/admin")
		{
			admin.Use(adminMiddleware)
			admin.GET("/users", adminHandler.ListUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.PATCH("/users/:id/position", adminHandler.UpdatePosition)
			admin.GET("/users/:id/avatar", adminHandler.GetUserAvatar)
			admin.PATCH("/users/:id/role", adminHandler.UpdateRole)
			admin.PATCH("/users/:id/manager", adminHandler.AssignManager)
			admin.DELETE("/users/:id/manager", adminHandler.RemoveManager)
			admin.GET("/managers/:id/employees", adminHandler.ListEmployeesByManager)

		}
		manager := authorized.Group("/manager")
		{
			manager.Use(managerMiddleware)
			manager.GET("/employees", userHandler.GetEmployeesByManager)
			plans := manager.Group("/plans")
			{
				plans.GET("", planHandler.ListByManager)
				plans.POST("", planHandler.Create)
				plans.POST("/ai", planHandler.CreateAI)
				plans.GET("/:plan_id", planHandler.GetByID)
				plans.DELETE("/:plan_id", planHandler.Delete)
				plans.PATCH("/:plan_id/archive", planHandler.Archive)
				plans.PATCH("/:plan_id", planHandler.Update)
				plans.POST("/:plan_id/tasks", taskHandler.Create)
			}

			tasks := manager.Group("/tasks")
			{
				tasks.GET("/:task_id", taskHandler.GetByID)
				tasks.DELETE("/:task_id", taskHandler.Delete)
				tasks.PATCH("/:task_id", taskHandler.Update)
			}
		}
		employee := authorized.Group("/employee")
		{
			employee.Use(employeeMiddleware)
			skills := employee.Group("/skills")
			{
				skills.GET("", skillHandler.EmployeeList)
			}
			plans := employee.Group("/plans")
			{
				plans.GET("/:plan_id/test", testHandler.GetForEmployee)
				plans.POST("/:plan_id/test", testHandler.Submit)
				plans.GET("", planHandler.EmployeeGetPlans)
				plans.GET("/:plan_id", planHandler.EmployeeGetPlan)
			}
			tasks := employee.Group("/tasks")
			{
				tasks.PATCH("/:task_id/status", taskHandler.UpdateStatus)
			}
		}
	}

	return router
}
