package router

import (
	"workHub/internal/controller"
	"workHub/internal/repository"
	"workHub/internal/service"

	"crypto/rsa"
	"gorm.io/gorm"
)

type Dependencies struct {
	PublicKey            *rsa.PublicKey
	AuthController       *controller.AuthController
	BookController       *controller.BookController
	BorrowController     *controller.BorrowController
	RoleController       *controller.RoleController
	PermissionController *controller.PermissionController
	WorkspaceController  *controller.WorkspaceController
	PermissionRepo       repository.PermissionRepository
	AuthRepo             repository.AuthRepository
}

func InitDependencies(db *gorm.DB) (*Dependencies, error) {
	// Khởi tạo JWT config
	jwtConfig, err := InitJWTConfig()
	if err != nil {
		return nil, err
	}

	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	sessionRedisRepo := repository.NewSessionRedisRepository()
	authService := service.NewAuthService(authRepo, sessionRepo, sessionRedisRepo, jwtConfig.Service)
	authController := controller.NewAuthController(authService)

	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleController := controller.NewRoleController(roleService)

	permissionRepo := repository.NewPermissionRepository(db)
	permissionService := service.NewPermissionService(permissionRepo, roleRepo)
	permissionController := controller.NewPermissionController(permissionService)

	workspaceRepo := repository.NewWorkspaceRepository(db)
	workspaceService := service.NewWorkspaceService(workspaceRepo, authRepo)
	workspaceController := controller.NewWorkspaceController(workspaceService)

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookController := controller.NewBookController(bookService)

	borrowRepo := repository.NewBorrowRepository(db)
	borrowService := service.NewBorrowService(borrowRepo, bookRepo, authRepo, db)
	borrowController := controller.NewBorrowController(borrowService)

	return &Dependencies{
		PublicKey:            jwtConfig.PublicKey,
		AuthController:       authController,
		BookController:       bookController,
		BorrowController:     borrowController,
		RoleController:       roleController,
		PermissionController: permissionController,
		WorkspaceController:  workspaceController,
		PermissionRepo:       permissionRepo,
		AuthRepo:             authRepo,
	}, nil
}

