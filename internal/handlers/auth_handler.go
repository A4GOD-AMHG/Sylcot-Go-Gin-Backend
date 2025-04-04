package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/A4GOD-AMHG/sylcot-go-gin-backend/docs"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/models"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/repositories"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo repositories.AuthRepository
}

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"Password*1"`
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Registration data"
// @Success 201 {object} object{message=string}
// @Failure 400 {object} object{error=string,details=object}
// @Failure 409 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/v1/auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": map[string]interface{}{},
		})
		return
	}

	log.Println(user)

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

	verificationLink := os.Getenv("FRONT_URL") + "/auth/verify-email?token=" + newUser.Token
	if err := utils.SendVerificationEmail(user.Email, verificationLink); err != nil {
		log.Printf("Could not send verification email to %s: %v", user.Email, err)
	}

	fmt.Println(verificationLink)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully. Please verify your email."})
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"Password*1"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token and user info
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} object{token=string,user=models.UserDTO}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/v1/auth/login [post]
// Ejemplo de request:
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

	userDTO := user.ToDTO()

	c.JSON(http.StatusOK, gin.H{"token": jwtToken, "user": userDTO})
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Validate email verification token
// @Tags authentication
// @Produce json
// @Param token query string true "Verification token from email"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/v1/auth/verify-email [get]
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

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset instructions to email
// @Tags authentication
// @Accept json
// @Produce json
// @Param email body object{email=string} true "Registered email address"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/v1/auth/forgot-password [post]
func (ah *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	user, err := ah.Repo.FindByEmail(req.Email)
	if err != nil {

		c.JSON(http.StatusOK, gin.H{"message": "If an account exists, a reset link has been sent"})
		return
	}

	resetToken := uuid.NewString()
	user.ResetToken = resetToken
	if err := ah.Repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate reset token"})
		return
	}

	resetLink := os.Getenv("API_URL") + "/api/v1/auth/reset-password?token=" + resetToken
	if err := utils.SendResetPasswordEmail(user.Email, resetLink); err != nil {
		log.Printf("Error sending reset email to %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If an account exists, a reset link has been sent"})
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,password" example:"Password*1"`
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Set new password using reset token
// @Tags authentication
// @Accept json
// @Produce json
// @Param resetData body ResetPasswordRequest true "Reset password data"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/v1/auth/reset-password [post]
func (ah *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8,password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user, err := ah.Repo.FindByResetToken(req.Token)
	if err != nil {
		if errors.Is(err, repositories.ErrTokenNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	if err := ah.Repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Password successfully updated for %s", user.Email)})
}
