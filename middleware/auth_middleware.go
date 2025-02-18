package middleware

import (
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/alastor-4/sylcot-go-gin-backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Email  string `json:"email"`
	UserID int    `json:"userId"`
	jwt.RegisteredClaims
}

// AuthMiddleware godoc
// @Security ApiKeyAuth
// @Description JWT Authentication Middleware with refresh detection
// @Param Authorization header string true "JWT Token" default(Bearer <token>)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		claims := &CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims.ExpiresAt != nil {
			remaining := time.Until(claims.ExpiresAt.Time)
			if remaining < 5*time.Minute {
				newTokenString, err := controllers.GenerateJWT(claims.Email, claims.UserID)
				if err == nil {
					c.Header("X-Refresh-Token", newTokenString)
				}
			}
		}

		c.Set("userEmail", claims.Email)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
