package models

import (
	"unicode"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// User represents application user
// @SWG.Definition(
//
//	required: ["name", "email", "password"],
//	properties: {
//	    "id": {type: "integer", example: 1},
//	    "name": {type: "string", example: "John Doe", minLength: 2, maxLength: 50},
//	    "email": {type: "string", format: "email", example: "user@example.com"},
//	    "password": {type: "string", format: "password", example: "P@ssw0rd!", minLength: 8},
//	    "is_verified": {type: "boolean", example: false},
//	    "token": {type: "string", example: "550e8400-e29b-41d4-a716-446655440000"}
//	}
//
// )
type User struct {
	gorm.Model
	Name         string `gorm:"size:255" json:"name" validate:"required,min=2,max=50"`
	Email        string `gorm:"unique;size:255" json:"email" validate:"required,email"`
	Password     string `gorm:"size:255" json:"password" validate:"required,min=8,password"`
	IsVerified   bool   `gorm:"default:false" json:"is_verified"`
	RefreshToken string `gorm:"size:255" json:"refresh_token"`
	Token        string `gorm:"size:255" json:"token"`
}

func GetValidationMessages(err error) map[string][]string {
	errors := make(map[string][]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch field {
			case "Email":
				switch tag {
				case "required":
					errors["email"] = append(errors["email"], "Email is required")
				case "email":
					errors["email"] = append(errors["email"], "Invalid email format")
				}
			case "Password":
				switch tag {
				case "required":
					errors["password"] = append(errors["password"], "Password is required")
				case "min":
					errors["password"] = append(errors["password"], "Password must be at least 8 characters")
				case "password":
					errors["password"] = append(errors["password"], "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
				}
			case "Name":
				switch tag {
				case "required":
					errors["name"] = append(errors["name"], "Name is required")
				case "min", "max":
					errors["name"] = append(errors["name"], "Name must be between 2 and 50 characters")
				}
			}
		}
	}

	return errors
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (u *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", passwordValidator)
	return validate.Struct(u)
}

func MigrateUsers(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
