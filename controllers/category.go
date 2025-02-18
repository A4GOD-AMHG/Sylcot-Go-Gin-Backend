package controllers

import (
	"net/http"

	"github.com/alastor-4/sylcot-go-gin-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

// GetCategories godoc
// @Summary Get categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Category
// @Failure 500 {object} object{error=string}
// @Router /api/categories [get]
func (cc *CategoryController) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := cc.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
