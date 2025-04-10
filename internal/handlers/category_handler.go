package handlers

import (
	"net/http"

	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/models"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/repositories"
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
// @Success 200 {array} models.CategoryDTO
// @Failure 500 {object} object{error=string}
// @Router /api/categories [get]
func (ch *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := ch.repo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}

	var categoriesDTO []*models.CategoryDTO
	for _, category := range categories {
		categoriesDTO = append(categoriesDTO, category.ToDTO())
	}

	c.JSON(http.StatusOK, categoriesDTO)
}
