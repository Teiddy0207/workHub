package router

import (
	"workHub/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	
	// CORS middleware - phải đặt trước tất cả routes
	r.Use(middleware.CORSMiddleware())
	
	deps, err := InitDependencies(db)
	if err != nil {
		panic(err)
	}
	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", deps.AuthController.Register)
		auth.POST("/login", deps.AuthController.Login)
		auth.GET("/me", middleware.AuthMiddleware(deps.PublicKey), middleware.RoleMiddleware(deps.AuthRepo), deps.AuthController.GetMe)
	}

	// Public book routes
	books := r.Group("/api/books")
	{
		books.GET("", deps.BookController.List)
		books.GET("/:id", deps.BookController.GetByID)
		books.GET("/:id/borrows", deps.BorrowController.GetByBookID)
	}

	protected := r.Group("/api/v1/")
	protected.Use(middleware.AuthMiddleware(deps.PublicKey))
	protected.Use(middleware.RoleMiddleware(deps.AuthRepo))
	{
		users := protected.Group("/users")
		{
			users.GET("", deps.AuthController.GetListUser)
			users.GET("/:id", deps.AuthController.GetUserByID)
			users.PUT("/:id", deps.AuthController.UpdateUser)
			users.DELETE("/:id", deps.AuthController.DeleteUser)
			users.GET("/:id/borrows", deps.BorrowController.GetByUserID)

			users.POST("/:id/roles", deps.PermissionController.AssignRolesToUser)
			users.DELETE("/:id/roles", deps.PermissionController.RemoveRolesFromUser)
			users.GET("/:id/permissions", deps.PermissionController.GetUserPermissions)
		}

		books := protected.Group("/books")
		{
			books.POST("", deps.BookController.Create)
			books.PUT("/:id", deps.BookController.Update)
			books.DELETE("/:id", deps.BookController.Delete)
			books.GET("/search", deps.BookController.Search)
		}

		borrows := protected.Group("/borrows")
		{
			borrows.POST("", deps.BorrowController.BorrowBook)
			borrows.PUT("/:id/return", deps.BorrowController.ReturnBook)
			borrows.GET("/active", deps.BorrowController.GetActiveBorrows)
		}

		roles := protected.Group("/roles")
		{
			roles.POST("", deps.RoleController.CreateRole)
			roles.GET("", deps.RoleController.ListRoles)
			roles.GET("/:id", deps.RoleController.GetRoleByID)
			roles.PUT("/:id", deps.RoleController.UpdateRole)
			roles.DELETE("/:id", deps.RoleController.DeleteRole)

			roles.POST("/:id/permissions", deps.PermissionController.AssignPermissionsToRole)
			roles.DELETE("/:id/permissions", deps.PermissionController.RemovePermissionsFromRole)
			roles.GET("/:id/permissions", deps.PermissionController.GetRoleWithPermissions)
		}

		permissions := protected.Group("/permissions")
		{
			permissions.POST("", deps.PermissionController.CreatePermission)
			permissions.GET("", deps.PermissionController.ListPermissions)
			permissions.GET("/:id", deps.PermissionController.GetPermissionByID)
			permissions.PUT("/:id", deps.PermissionController.UpdatePermission)
			permissions.DELETE("/:id", deps.PermissionController.DeletePermission)
		}

		workspaces := protected.Group("/workspaces")
		{
			workspaces.POST("", deps.WorkspaceController.CreateWorkspace)
			workspaces.GET("", deps.WorkspaceController.ListWorkspaces)
			workspaces.GET("/:id", deps.WorkspaceController.GetWorkspaceByID)
			workspaces.PUT("/:id", deps.WorkspaceController.UpdateWorkspace)
			workspaces.DELETE("/:id", deps.WorkspaceController.DeleteWorkspace)
		}
	}

	return r
}
