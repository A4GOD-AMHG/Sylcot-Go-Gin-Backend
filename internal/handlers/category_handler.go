package handlers

import (
	"net/http"

	_ "github.com/alastor-4/sylcot-go-gin-backend/internal/models"
	"github.com/alastor-4/sylcot-go-gin-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo repositories.CategoryRepository
}

func NewCategoryHandler(repo repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
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
func (ch *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := ch.repo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
