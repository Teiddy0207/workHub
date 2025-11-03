package middleware

import (
	"context"
	"fmt"
	"workHub/internal/repository"
	"workHub/logger"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware kiểm tra xem user có permission cụ thể không
// permissionCode là code của permission cần kiểm tra (ví dụ: "user.create", "role.delete")
func PermissionMiddleware(permissionRepo repository.PermissionRepository, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy user_id từ context (đã được set bởi AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Warn("middleware", "PermissionMiddleware", "User ID not found in context")
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    401,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			logger.Warn("middleware", "PermissionMiddleware", "Invalid user ID type")
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    401,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Kiểm tra permission
		hasPermission, err := permissionRepo.UserHasPermission(context.Background(), userIDStr, permissionCode)
		if err != nil {
			logger.Error("middleware", "PermissionMiddleware", fmt.Sprintf("Failed to check permission: %v", err))
			c.JSON(500, gin.H{
				"status":  "error",
				"code":    500,
				"message": "Internal server error",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn("middleware", "PermissionMiddleware", fmt.Sprintf("User %s does not have permission: %s", userIDStr, permissionCode))
			c.JSON(403, gin.H{
				"status":  "error",
				"code":    403,
				"message": fmt.Sprintf("Forbidden: You don't have permission '%s'", permissionCode),
			})
			c.Abort()
			return
		}

		logger.Info("middleware", "PermissionMiddleware", fmt.Sprintf("User %s has permission: %s", userIDStr, permissionCode))
		c.Next()
	}
}


