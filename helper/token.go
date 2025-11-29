package helper

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrTokenNotFound = errors.New("token not found in header")
	ErrInvalidToken  = errors.New("invalid token format")
)

func GenerateResetToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"type":  "reset_password",
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// GetTokenFromHeader lấy JWT token từ Authorization header
// Format: "Bearer <token>"
func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrTokenNotFound
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidToken
	}

	return parts[1], nil
}

// GetUserEmailFromToken lấy email từ JWT token
// Token đã được parse bởi AuthMiddleware và lưu email vào context
func GetUserEmailFromToken(c *gin.Context) (string, error) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", errors.New("user email not found in context, token may not be authenticated")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.New("invalid user email type in context")
	}

	return emailStr, nil
}
