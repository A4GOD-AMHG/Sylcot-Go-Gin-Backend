package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func getJWTExpiration() time.Duration {
	minutesStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	if minutesStr == "" {
		return time.Minute * 4320
	}
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return time.Minute * 4320
	}
	return time.Minute * time.Duration(minutes)
}

func GenerateJWT(email string, id int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expiration := getJWTExpiration()
	claims := jwt.MapClaims{
		"email":  email,
		"userId": id,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
