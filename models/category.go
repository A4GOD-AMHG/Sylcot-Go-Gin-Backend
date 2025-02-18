package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Title    string `gorm:"unique;size:50" json:"title"`
	Color    string `gorm:"unique;size:50" json:"color"`
	IconName string `gorm:"unique;size:50" json:"icon_name"`
	Tasks    []Task `gorm:"foreignKey:CategoryID" json:"-"`
}

func MigrateCategories(db *gorm.DB) error {
	return db.AutoMigrate(&Category{})
}

func (c *Category) Setup(db *gorm.DB) error {
	categories := []Category{
		{Title: "Home", Color: "#FFAB91", IconName: "home"},
		{Title: "Work", Color: "#80D8FF", IconName: "work_outline"},
		{Title: "Sports", Color: "#81C784", IconName: "sports_soccer"},
		{Title: "Couple", Color: "#F48FB1", IconName: "favorite"},
		{Title: "Health", Color: "#EF9A9A", IconName: "favorite_border"},
		{Title: "Study", Color: "#FFF59D", IconName: "menu_book"},
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
