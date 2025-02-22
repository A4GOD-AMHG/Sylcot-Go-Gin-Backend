package routes

import (
	"github.com/alastor-4/sylcot-go-gin-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupTaskRoutes(router *gin.RouterGroup, handler *handlers.TaskHandler) {
	router.GET("/tasks", handler.GetTasks)
	router.POST("/tasks", handler.CreateTask)
	router.PUT("/tasks/:id", handler.UpdateTask)
	router.DELETE("/tasks/:id", handler.DeleteTask)
	router.PATCH("/tasks/:id/complete", handler.ToggleTask)
}
