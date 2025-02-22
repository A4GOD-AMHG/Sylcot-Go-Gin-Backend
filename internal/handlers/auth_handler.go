package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/alastor-4/sylcot-go-gin-backend/docs"
	"github.com/alastor-4/sylcot-go-gin-backend/internal/models"
	"github.com/alastor-4/sylcot-go-gin-backend/internal/repositories"
	"github.com/alastor-4/sylcot-go-gin-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo repositories.AuthRepository
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body models.User true "Registration data"
// @Success 201 {object} map[string]interface{} "message: User registered successfully..."
// @Failure 400 {object} map[string]interface{} "error: Validation failed, details: field errors"
// @Failure 409 {object} map[string]interface{} "error: User already exists"
// @Failure 500 {object} map[string]interface{} "error: Internal server error"
// @Router /auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": map[string]interface{}{},
		})
		return
	}

	if err := user.Validate(); err != nil {
		validationErrors := models.GetValidationMessages(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	_, err := ah.Repo.FindByEmail(user.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with that email already registered"})
		return
	} else if !errors.Is(err, repositories.ErrUserNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}

	newUser := models.User{
		Name:       user.Name,
		Email:      user.Email,
		Password:   string(hashedPassword),
		IsVerified: false,
		Token:      uuid.NewString(),
	}

	if err := ah.Repo.CreateUser(&newUser); err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register the user"})
		}
		return
	}

	verificationLink := "http://localhost:8080/auth/verify-email?token=" + newUser.Token
	if err := utils.SendVerificationEmail(user.Email, verificationLink); err != nil {
		log.Printf("Could not send verification email to %s: %v", user.Email, err)
	}

	fmt.Println(verificationLink)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully. Please verify your email."})
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "token: JWT string"
// @Failure 400 {object} map[string]interface{} "error: Invalid data"
// @Failure 401 {object} map[string]interface{} "error: Invalid credentials"
// @Failure 403 {object} map[string]interface{} "error: Email not verified"
// @Failure 500 {object} map[string]interface{} "error: Internal server error"
func (ah *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user, err := ah.Repo.FindByEmail(loginData.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email first"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	jwtToken, err := utils.GenerateJWT(user.Email, int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT Token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Validate email verification token
// @Tags authentication
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{} "message: Verification success message"
// @Failure 400 {object} map[string]interface{} "error: Token required"
// @Failure 404 {object} map[string]interface{} "error: Invalid token"
// @Failure 500 {object} map[string]interface{} "error: Internal server error"
// @Router /auth/verify-email [get]
func (ah *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	user, err := ah.Repo.FindByToken(token)
	if err != nil {
		if errors.Is(err, repositories.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	user.IsVerified = true
	user.Token = ""

	if err := ah.Repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User with email %s verified successfully", user.Email),
	})
}
