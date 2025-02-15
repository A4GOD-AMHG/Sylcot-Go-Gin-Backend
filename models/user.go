package models

import (
	"unicode"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"size:255" json:"name" validate:"required,min=2,max=50"`
	Email      string `gorm:"unique;size:255" json:"email" validate:"required,email"`
	Image      string `gorm:"size:255" json:"image" validate:"omitempty,url"`
	Password   string `gorm:"size:255" json:"password" validate:"required,min=8,password"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`
	Token      string `gorm:"size:255" json:"token"`
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
