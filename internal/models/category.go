package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a task category
// @SWG.Definition(
//
//	required: ["title", "color", "icon_name"],
//	properties: {
//	    "id": {type: "integer", example: 1},
//	    "created_at": {type: "string", format: "date-time", example: "2025-03-27T14:45:46Z"},
//	    "updated_at": {type: "string", format: "date-time", example: "2025-03-27T14:45:46Z"},
//	    "deleted_at": {type: "string", format: "date-time", example: "null", x-nullable: true},
//	    "title": {type: "string", example: "Work", maxLength: 50},
//	    "color": {type: "string", example: "#80D8FF", maxLength: 50},
//	    "icon_name": {type: "string", example: "work_outline", maxLength: 50}
//	}
//
// )
type Category struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
	Title     string     `gorm:"unique;size:50" json:"title"`
	Color     string     `gorm:"unique;size:50" json:"color"`
	IconName  string     `gorm:"unique;size:50" json:"icon_name"`
	Tasks     []Task     `gorm:"foreignKey:CategoryID" json:"-"`
}

func MigrateCategories(db *gorm.DB) error {
	return db.AutoMigrate(&Category{})
}

func (c *Category) Setup(db *gorm.DB) error {
	categories := []Category{
		{Title: "Home", Color: "#FFAB91", IconName: "home"},
		{Title: "Work", Color: "#80D8FF", IconName: "briefcase"},
		{Title: "Sports", Color: "#81C784", IconName: "soccer"},
		{Title: "Couple", Color: "#F48FB1", IconName: "heart"},
		{Title: "Health", Color: "#EF9A9A", IconName: "heart-outline"},
		{Title: "Study", Color: "#FFF59D", IconName: "book"},
		{Title: "Shopping", Color: "#CE93D8", IconName: "shopping_cart"},
		{Title: "Finance", Color: "#A5D6A7", IconName: "attach_money"},
		{Title: "Travel", Color: "#90CAF9", IconName: "flight"},
		{Title: "Social", Color: "#FFCC80", IconName: "groups"},
		{Title: "Creativity", Color: "#B39DDB", IconName: "palette"},
		{Title: "Pets", Color: "#BCAAA4", IconName: "pets"},
		{Title: "Meals", Color: "#FF8A65", IconName: "restaurant"},
		{Title: "Others", Color: "#B0BEC5", IconName: "more_horiz"},
	}

	for _, categorie := range categories {
		result := db.FirstOrCreate(&categorie, Category{Title: categorie.Title})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
