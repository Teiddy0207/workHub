package helper

import (
	"errors"
	"workHub/internal/dto"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
)

// GetUserID lấy user ID từ context (được set bởi AuthMiddleware)
func GetUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, ErrUserNotAuthenticated
	}
	if id, ok := userID.(int); ok {
		return id, nil
	}
	return 0, ErrUserNotAuthenticated
}

// GetUserEmail lấy user email từ context
func GetUserEmail(c *gin.Context) (string, error) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", ErrUserNotAuthenticated
	}
	if e, ok := email.(string); ok {
		return e, nil
	}
	return "", ErrUserNotAuthenticated
}

// GetUserUsername lấy username từ context
func GetUserUsername(c *gin.Context) (string, error) {
	username, exists := c.Get("user_username")
	if !exists {
		return "", ErrUserNotAuthenticated
	}
	if u, ok := username.(string); ok {
		return u, nil
	}
	return "", ErrUserNotAuthenticated
}

// GetUserInfo lấy toàn bộ user info từ context
func GetUserInfo(c *gin.Context) (dto.Users, error) {
	userInfo, exists := c.Get("user_info")
	if !exists {
		return dto.Users{}, ErrUserNotAuthenticated
	}
	if u, ok := userInfo.(dto.Users); ok {
		return u, nil
	}
	return dto.Users{}, ErrUserNotAuthenticated
}

