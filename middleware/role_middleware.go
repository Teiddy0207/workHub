package middleware

import (
	"context"
	"fmt"
	"workHub/internal/repository"
	"workHub/logger"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware load user role từ DB và lưu vào context
func RoleMiddleware(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.Next()
			return
		}

		// Load user từ DB để lấy role
		user, err := authRepo.GetUserByID(context.Background(), userIDStr)
		if err != nil {
			logger.Warn("middleware", "RoleMiddleware", fmt.Sprintf("Failed to load user role for userID=%s: %v", userIDStr, err))
			c.Next()
			return
		}

		// Lưu role vào context
		c.Set("user_role", user.Role)
		logger.Info("middleware", "RoleMiddleware", fmt.Sprintf("User role loaded: userID=%s, role=%s", userIDStr, user.Role))
		c.Next()
	}
}

