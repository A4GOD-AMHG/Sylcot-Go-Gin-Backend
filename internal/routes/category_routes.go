package routes

import (
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.RouterGroup, handler *handlers.CategoryHandler) {
	router.GET("/categories", handler.GetCategories)
}
