package router

import (
	"fmt"
	internalconfig "workHub/internal/config"
	"workHub/internal/controller"
	"workHub/internal/repository"
	"workHub/internal/service"
	"workHub/logger"
	"workHub/middleware"
	"workHub/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// db được truyền vào từ main, không cần connect lại

	// Load config để lấy JWT config
	cfg, err := internalconfig.LoadConfig()
	if err != nil {
		logger.Error("config", "LoadConfig", fmt.Sprintf("Config load failed: %v", err))
		panic(err)
	}

	// Khởi tạo JWT service
	jwtService, err := service.NewJWTService(*cfg)
	if err != nil {
		logger.Error("service", "NewJWTService", fmt.Sprintf("JWT service init failed: %v", err))
		panic(err)
	}

	// Parse public key để dùng trong middleware
	_, publicKey, _, err := jwt.ParseKey(cfg.Jwt)
	if err != nil {
		logger.Error("router", "InitRouter", fmt.Sprintf("Failed to parse JWT keys: %v", err))
		panic(err)
	}

	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	authService := service.NewAuthService(authRepo, sessionRepo, jwtService)
	authController := controller.NewAuthController(authService)

	// Role services
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleController := controller.NewRoleController(roleService)

	// Public routes (không cần authentication)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		// auth.POST("/register", authController.Register)
	}

	// Protected routes (cần authentication)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(publicKey))
	{
		// Auth routes cần đăng nhập
		protected.GET("/auth/users", authController.GetListUser)

		// Role routes cần đăng nhập
		roles := protected.Group("/roles")
		{
			roles.POST("", roleController.CreateRole)
			roles.GET("", roleController.ListRoles)
			roles.GET("/:id", roleController.GetRoleByID)
			roles.PUT("/:id", roleController.UpdateRole)
			roles.DELETE("/:id", roleController.DeleteRole)
		}
	}

	return r
}
