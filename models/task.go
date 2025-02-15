package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Priority string

const (
	High   Priority = "high"
	Medium Priority = "medium"
	Low    Priority = "low"
)

type Task struct {
	gorm.Model
	Title      string   `gorm:"size:100;not null" json:"title" validate:"required,min=3,max=100"`
	Priority   Priority `gorm:"type:varchar(10);default:'medium'" json:"priority"`
	Status     bool     `gorm:"default:false" json:"status"`
	CategoryID uint     `gorm:"not null" json:"category_id" validate:"required"`
	UserID     uint     `gorm:"not null" json:"user_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
	User       User     `gorm:"foreignKey:UserID" json:"-"`
}

type TaskRequest struct {
	Title      string   `json:"title" validate:"required,min=3,max=100"`
	Priority   Priority `json:"priority" validate:"oneof=high medium low"`
	CategoryID uint     `json:"category_id" validate:"required"`
}

func ValidateTaskRequest(taskReq TaskRequest) error {
	validate := validator.New()
	validate.RegisterValidation("priority", validatePriority)
	return validate.Struct(taskReq)
}

func validatePriority(fl validator.FieldLevel) bool {
	priority := fl.Field().Interface().(Priority)
	switch priority {
	case High, Medium, Low:
		return true
	default:
		return false
	}
}

func GetTaskValidationMessages(err error) map[string][]string {
	errors := make(map[string][]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch field {
			case "Title":
				switch tag {
				case "required":
					errors["title"] = append(errors["title"], "Title is required")
				case "min", "max":
					errors["title"] = append(errors["title"], "Title must be between 3 and 100 characters")
				}
			case "Priority":
				errors["priority"] = append(errors["priority"], "Priority must be one of: high, medium, low")
			case "CategoryID":
				errors["category_id"] = append(errors["category_id"], "Category is required")
			}
		}
	}

	return errors
}

func MigrateTasks(db *gorm.DB) error {
	return db.AutoMigrate(&Task{})
}
