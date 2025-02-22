package routes

import (
	"github.com/alastor-4/sylcot-go-gin-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.RouterGroup, handler *handlers.CategoryHandler) {
	router.GET("/categories", handler.GetCategories)
}
