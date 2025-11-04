package router

import (
	"workHub/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	deps, err := InitDependencies(db)
	if err != nil {
		panic(err)
	}
	auth := r.Group("/auth")
	{
		auth.POST("/login", deps.AuthController.Login)
	}

	protected := r.Group("/api/v1/")
	protected.Use(middleware.AuthMiddleware(deps.PublicKey))
	{
		users := protected.Group("/users")
		{
			users.GET("", middleware.PermissionMiddleware(deps.PermissionRepo, "user.read"), deps.AuthController.GetListUser)

			users.POST("/:id/roles", deps.PermissionController.AssignRolesToUser)
			users.DELETE("/:id/roles", deps.PermissionController.RemoveRolesFromUser)
			users.GET("/:id/permissions", deps.PermissionController.GetUserPermissions)
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
	}

	return r
}
