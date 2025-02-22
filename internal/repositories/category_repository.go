package repositories

import (
	"errors"
	"fmt"

	"github.com/alastor-4/sylcot-go-gin-backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepository interface {
	GetAllCategories() ([]models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (cr *categoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := cr.db.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	return categories, nil
}
