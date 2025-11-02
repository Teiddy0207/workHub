package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"
	"workHub/pkg/jwt"
	"workHub/constant"
	"workHub/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware xác thực JWT token và lưu user info vào context
// publicKey được truyền vào để tránh load config mỗi request
func AuthMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("middleware", "AuthMiddleware", "Missing Authorization header")
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    401,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("middleware", "AuthMiddleware", "Invalid Authorization header format")
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    401,
				"message": "Invalid Authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Verify token
		claims, err := jwt.VerifyToken(context.Background(), publicKey, tokenStr)
		if err != nil {
			logger.Warn("middleware", "AuthMiddleware", fmt.Sprintf("Token verification failed: %v", err))
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    401,
				"message": constant.ErrUnAuthentication.Error(),
			})
			c.Abort()
			return
		}

		// Kiểm tra token type phải là access token
		if claims.Type != jwt.TokenTypeAccessToken {
			logger.Warn("middleware", "AuthMiddleware", "Token is not an access token")
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    401,
				"message": "Invalid token type",
			})
			c.Abort()
			return
		}

		// Lưu user info vào context
		c.Set("user_id", claims.UserInfo.ID)
		c.Set("user_email", claims.UserInfo.Email)
		c.Set("user_username", claims.UserInfo.Username)
		c.Set("user_info", claims.UserInfo)

		logger.Info("middleware", "AuthMiddleware", fmt.Sprintf("User authenticated: %s", claims.UserInfo.Email))

		c.Next()
	}
}

