package router

import (
	"workHub/config"
	internalconfig "workHub/internal/config"
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
		panic(err)
	}

	// Load config để lấy JWT config
	cfg, err := internalconfig.LoadConfig()
	if err != nil {
		logger.Error("config", "LoadConfig", "Config load failed", zap.Error(err))
		panic(err)
	}

	// Khởi tạo JWT service
	jwtService, err := service.NewJWTService(*cfg)
	if err != nil {
		logger.Error("service", "NewJWTService", "JWT service init failed", zap.Error(err))
		panic(err)
	}

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService)
	authController := controller.NewAuthController(authService)

	// Role services
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleController := controller.NewRoleController(roleService)

	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		// auth.POST("/register", authController.Register)
		auth.GET("/users", authController.GetListUser)
	}

	roles := r.Group("/roles")
	{
		roles.POST("", roleController.CreateRole)
		roles.GET("", roleController.ListRoles)
		roles.GET("/:id", roleController.GetRoleByID)
		roles.PUT("/:id", roleController.UpdateRole)
		roles.DELETE("/:id", roleController.DeleteRole)
	}

	return r
}
