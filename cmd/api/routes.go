package api

import (
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/handlers"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine,
	authHandler *handlers.AuthHandler,
	taskHandler *handlers.TaskHandler,
	categoryHandler *handlers.CategoryHandler) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.GET("/verify-email", authHandler.VerifyEmail)
	}

	api := router.Group("/api/v1").Use(middleware.AuthMiddleware())
	{
		api.GET("/tasks", taskHandler.GetTasks)
		api.POST("/tasks", taskHandler.CreateTask)
		api.PUT("/tasks/:id", taskHandler.UpdateTask)
		api.DELETE("/tasks/:id", taskHandler.DeleteTask)
		api.PATCH("/tasks/:id/complete", taskHandler.ToggleTask)

		api.GET("/categories", categoryHandler.GetCategories)
	}
}
