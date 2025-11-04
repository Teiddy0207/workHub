package router

import (
	"workHub/internal/controller"
	"workHub/internal/repository"
	"workHub/internal/service"

	"crypto/rsa"
	"gorm.io/gorm"
)

type Dependencies struct {
	PublicKey          *rsa.PublicKey
	AuthController     *controller.AuthController
	RoleController     *controller.RoleController
	PermissionController *controller.PermissionController
	PermissionRepo     repository.PermissionRepository
}

func InitDependencies(db *gorm.DB) (*Dependencies, error) {
	// Khởi tạo JWT config
	jwtConfig, err := InitJWTConfig()
	if err != nil {
		return nil, err
	}

	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	authService := service.NewAuthService(authRepo, sessionRepo, jwtConfig.Service)
	authController := controller.NewAuthController(authService)

	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleController := controller.NewRoleController(roleService)

	permissionRepo := repository.NewPermissionRepository(db)
	permissionService := service.NewPermissionService(permissionRepo, roleRepo)
	permissionController := controller.NewPermissionController(permissionService)

	return &Dependencies{
		PublicKey:           jwtConfig.PublicKey,
		AuthController:      authController,
		RoleController:      roleController,
		PermissionController: permissionController,
		PermissionRepo:       permissionRepo,
	}, nil
}

