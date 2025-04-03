package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Priority string

const (
	High   Priority = "high"
	Medium Priority = "medium"
	Low    Priority = "low"
)

// Task represents a user task
// @SWG.Definition(
//
//	required: ["title", "priority", "category_id", "user_id"],
//	properties: {
//	    "id": {type: "integer", example: 1},
//	    "created_at": {type: "string", format: "date-time", example: "2025-03-27T14:45:46Z"},
//	    "updated_at": {type: "string", format: "date-time", example: "2025-03-27T14:45:46Z"},
//	    "deleted_at": {type: "string", format: "date-time", example: "null", x-nullable: true},
//	    "title": {type: "string", example: "Complete project report", minLength: 3, maxLength: 255},
//	    "priority": {type: "string", enum: ["high", "medium", "low"], example: "medium"},
//	    "status": {type: "boolean", example: false},
//	    "category_id": {type: "integer", example: 2},
//	    "user_id": {type: "integer", example: 1},
//	    "category": {"$ref": "#/definitions/Category"},
//	    "user": {"$ref": "#/definitions/User"}
//	}
//
// )
type Task struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
	Title      string     `gorm:"size:255;not null;uniqueIndex:idx_user_title" json:"title" validate:"required,min=3,max=255"`
	Priority   Priority   `gorm:"type:varchar(10);default:'medium'" json:"priority"`
	Status     bool       `gorm:"default:false" json:"status"`
	CategoryID uint       `gorm:"not null" json:"category_id" validate:"required"`
	UserID     uint       `gorm:"not null;uniqueIndex:idx_user_title" json:"user_id"`
	Category   Category   `gorm:"foreignKey:CategoryID" json:"category"`
	User       User       `gorm:"foreignKey:UserID" json:"user"`
}

// TaskRequest represents the payload for creating/updating a task
// @SWG.Definition(
//
//	required: ["title", "category_id"],
//	properties: {
//	    "title": {type: "string", example: "Buy groceries", minLength: 3, maxLength: 100},
//	    "priority": {type: "string", enum: ["high", "medium", "low"], example: "medium"},
//	    "category_id": {type: "integer", example: 3}
//	}
//
// )
type TaskRequest struct {
	Title      string   `json:"title" validate:"required,min=3,max=100"`
	Priority   Priority `json:"priority" validate:"oneof=high medium low"`
	CategoryID uint     `json:"category_id" validate:"required"`
}

type TaskDTO struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	Title     string      `json:"title"`
	Priority  Priority    `json:"priority"`
	Status    bool        `json:"status"`
	Category  CategoryDTO `json:"category"`
}

func (t *Task) ToDTO() *TaskDTO {
	return &TaskDTO{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		Title:     t.Title,
		Priority:  t.Priority,
		Status:    t.Status,
		Category:  *t.Category.ToDTO(),
	}
}

func ValidateTaskRequest(taskReq TaskRequest) error {
	validate := validator.New()
	validate.RegisterValidation("priority", validatePriority)
	return validate.Struct(taskReq)
}

func IsValidPriority(p Priority) bool {
	switch p {
	case High, Medium, Low:
		return true
	default:
		return false
	}
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

	err := db.Migrator().CreateIndex(&Task{}, "idx_user_title")
	if err != nil {
		return err
	}

	return db.AutoMigrate(&Task{})
}
