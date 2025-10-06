package router

import (
	"workHub/config"
	"workHub/internal/controller"
	"workHub/internal/repository"
	"workHub/internal/service"
	"workHub/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	db, err := config.ConnectDatabase()
	if err != nil {
		logger.Error("config", "ConnectDatabase", "DB connect failed", zap.Error(err))
	}

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
	}

	return r
}
