package routes

import (
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, handler *handlers.AuthHandler) {
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.GET("/verify-email", handler.VerifyEmail)
	router.POST("/forgot-password", handler.ForgotPassword)
	router.POST("/reset-password", handler.ResetPassword)
}
