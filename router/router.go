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

	cfg, err := internalconfig.LoadConfig()
	if err != nil {
		logger.Error("config", "LoadConfig", fmt.Sprintf("Config load failed: %v", err))
		panic(err)
	}

	jwtService, err := service.NewJWTService(*cfg)
	if err != nil {
		logger.Error("service", "NewJWTService", fmt.Sprintf("JWT service init failed: %v", err))
		panic(err)
	}
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

	// Permission services
	permissionRepo := repository.NewPermissionRepository(db)
	permissionService := service.NewPermissionService(permissionRepo, roleRepo)
	permissionController := controller.NewPermissionController(permissionService)

	// Public routes (không cần authentication)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		// auth.POST("/register", authController.Register)
	}

	// Protected routes (cần authentication)
	protected := r.Group("/api/v1/")
	protected.Use(middleware.AuthMiddleware(publicKey))
	{
		// User routes với permission check
		users := protected.Group("/users")
		{
			// Route xem danh sách users - cần quyền read
			users.GET("", middleware.PermissionMiddleware(permissionRepo, "user.read"), authController.GetListUser)
			
			// Route quản lý roles của user
			users.POST("/:id/roles", permissionController.AssignRolesToUser)
			users.DELETE("/:id/roles", permissionController.RemoveRolesFromUser)
			users.GET("/:id/permissions", permissionController.GetUserPermissions)
			
			// Ví dụ: Route xóa user - cần quyền delete (cần tạo method DeleteUser trong AuthController)
			// users.DELETE("/:id", middleware.PermissionMiddleware(permissionRepo, "user.delete"), authController.DeleteUser)
		}

		// Role routes cần đăng nhập
		roles := protected.Group("/roles")
		{
			roles.POST("", roleController.CreateRole)
			roles.GET("", roleController.ListRoles)
			roles.GET("/:id", roleController.GetRoleByID)
			roles.PUT("/:id", roleController.UpdateRole)
			roles.DELETE("/:id", roleController.DeleteRole)

			// Role-Permission management
			roles.POST("/:id/permissions", permissionController.AssignPermissionsToRole)
			roles.DELETE("/:id/permissions", permissionController.RemovePermissionsFromRole)
			roles.GET("/:id/permissions", permissionController.GetRoleWithPermissions)
		}

		permissions := protected.Group("/permissions")
		{
			permissions.POST("", permissionController.CreatePermission)
			permissions.GET("", permissionController.ListPermissions)
			permissions.GET("/:id", permissionController.GetPermissionByID)
			permissions.PUT("/:id", permissionController.UpdatePermission)
			permissions.DELETE("/:id", permissionController.DeletePermission)
		}
	}

	return r
}
