package repository

import (
	"context"
	"fmt"
	"workHub/internal/entity"
	"workHub/logger"
	"workHub/pkg/params"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

type PermissionRepository interface {
	CreatePermission(ctx context.Context, permission *entity.Permission) error
	GetPermissionByID(ctx context.Context, id string) (entity.Permission, error)
	GetPermissionByCode(ctx context.Context, code string) (entity.Permission, error)
	UpdatePermission(ctx context.Context, id string, permission *entity.Permission) error
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context, params params.QueryParams) (entity.PaginatedPermissions, error)
	
	// RolePermission methods
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error
	RemovePermissionsFromRole(ctx context.Context, roleID string, permissionIDs []string) error
	GetPermissionsByRoleID(ctx context.Context, roleID string, params params.QueryParams) (entity.PaginatedPermissions, error)
	
	// UserRole methods
	AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error
	RemoveRolesFromUser(ctx context.Context, userID string, roleIDs []string) error
	GetRolesByUserID(ctx context.Context, userID string) ([]entity.Role, error)
	
	// Permission check methods
	UserHasPermission(ctx context.Context, userID string, permissionCode string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]entity.Permission, error)
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	logger.Info("repository", "CreatePermission", fmt.Sprintf("Creating permission: %s (%s)", permission.Name, permission.Code))

	err := r.db.WithContext(ctx).Create(permission).Error
	if err != nil {
		logger.Error("repository", "CreatePermission", fmt.Sprintf("Failed to create permission: %v", err))
		return err
	}

	logger.Info("repository", "CreatePermission", fmt.Sprintf("Permission created successfully: %s", permission.ID))
	return nil
}

func (r *permissionRepository) GetPermissionByID(ctx context.Context, id string) (entity.Permission, error) {
	logger.Info("repository", "GetPermissionByID", fmt.Sprintf("Getting permission by ID: %s", id))

	var permission entity.Permission
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&permission).Error

	if err != nil {
		logger.Error("repository", "GetPermissionByID", fmt.Sprintf("Permission not found: %s, error: %v", id, err))
		return entity.Permission{}, err
	}

	logger.Info("repository", "GetPermissionByID", fmt.Sprintf("Permission found: %s", permission.Name))
	return permission, nil
}

func (r *permissionRepository) GetPermissionByCode(ctx context.Context, code string) (entity.Permission, error) {
	logger.Info("repository", "GetPermissionByCode", fmt.Sprintf("Getting permission by code: %s", code))

	var permission entity.Permission
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&permission).Error

	if err != nil {
		logger.Error("repository", "GetPermissionByCode", fmt.Sprintf("Permission not found: %s, error: %v", code, err))
		return entity.Permission{}, err
	}

	logger.Info("repository", "GetPermissionByCode", fmt.Sprintf("Permission found: %s", permission.Name))
	return permission, nil
}

func (r *permissionRepository) UpdatePermission(ctx context.Context, id string, permission *entity.Permission) error {
	logger.Info("repository", "UpdatePermission", fmt.Sprintf("Updating permission: %s", id))

	// Chỉ update các field không nil
	updates := make(map[string]interface{})
	if permission.Name != "" {
		updates["name"] = permission.Name
	}
	if permission.Code != "" {
		updates["code"] = permission.Code
	}
	if permission.Action != "" {
		updates["action"] = permission.Action
	}
	if permission.Description != "" {
		updates["description"] = permission.Description
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		logger.Error("repository", "UpdatePermission", fmt.Sprintf("Failed to update permission: %s, error: %v", id, err))
		return err
	}

	// Update is_active nếu có
	if permission.ID != "" {
		err = r.db.WithContext(ctx).
			Model(&entity.Permission{}).
			Where("id = ?", id).
			Update("is_active", permission.IsActive).Error
		if err != nil {
			logger.Error("repository", "UpdatePermission", fmt.Sprintf("Failed to update is_active: %v", err))
			return err
		}
	}

	logger.Info("repository", "UpdatePermission", fmt.Sprintf("Permission updated successfully: %s", id))
	return nil
}

func (r *permissionRepository) DeletePermission(ctx context.Context, id string) error {
	logger.Info("repository", "DeletePermission", fmt.Sprintf("Deleting permission: %s", id))

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.Permission{}).Error

	if err != nil {
		logger.Error("repository", "DeletePermission", fmt.Sprintf("Failed to delete permission: %s, error: %v", id, err))
		return err
	}

	logger.Info("repository", "DeletePermission", fmt.Sprintf("Permission deleted successfully: %s", id))
	return nil
}

func (r *permissionRepository) ListPermissions(ctx context.Context, params params.QueryParams) (entity.PaginatedPermissions, error) {
	logger.Info("repository", "ListPermissions", fmt.Sprintf("Listing permissions with params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	var permissions []entity.Permission
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Permission{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR action ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		logger.Error("repository", "ListPermissions", fmt.Sprintf("Failed to count permissions: %v", err))
		return entity.PaginatedPermissions{}, err
	}

	// Apply pagination
	offset := (params.PageNumber - 1) * params.PageSize
	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&permissions).Error; err != nil {
		logger.Error("repository", "ListPermissions", fmt.Sprintf("Failed to list permissions: %v", err))
		return entity.PaginatedPermissions{}, err
	}

	result := entity.PaginatedPermissions{
		Items:      permissions,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   params.PageSize,
	}

	logger.Info("repository", "ListPermissions", fmt.Sprintf("Listed %d permissions successfully", len(permissions)))
	return result, nil
}

// AssignPermissionsToRole gán nhiều permissions cho một role
func (r *permissionRepository) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	logger.Info("repository", "AssignPermissionsToRole", fmt.Sprintf("Assigning %d permissions to role: %s", len(permissionIDs), roleID))

	// Tạo các RolePermission records
	rolePermissions := make([]entity.RolePermission, len(permissionIDs))
	for i, permissionID := range permissionIDs {
		rolePermissions[i] = entity.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
	}

	// Sử dụng Create để thêm mới (GORM sẽ tự động bỏ qua duplicate nếu có unique constraint)
	err := r.db.WithContext(ctx).Create(&rolePermissions).Error
	if err != nil {
		logger.Error("repository", "AssignPermissionsToRole", fmt.Sprintf("Failed to assign permissions: %v", err))
		return err
	}

	logger.Info("repository", "AssignPermissionsToRole", "Permissions assigned successfully")
	return nil
}

// RemovePermissionsFromRole xóa permissions khỏi role
func (r *permissionRepository) RemovePermissionsFromRole(ctx context.Context, roleID string, permissionIDs []string) error {
	logger.Info("repository", "RemovePermissionsFromRole", fmt.Sprintf("Removing %d permissions from role: %s", len(permissionIDs), roleID))

	err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&entity.RolePermission{}).Error

	if err != nil {
		logger.Error("repository", "RemovePermissionsFromRole", fmt.Sprintf("Failed to remove permissions: %v", err))
		return err
	}

	logger.Info("repository", "RemovePermissionsFromRole", "Permissions removed successfully")
	return nil
}

// GetPermissionsByRoleID lấy permissions của một role có phân trang
func (r *permissionRepository) GetPermissionsByRoleID(ctx context.Context, roleID string, params params.QueryParams) (entity.PaginatedPermissions, error) {
	logger.Info("repository", "GetPermissionsByRoleID", fmt.Sprintf("Getting permissions for role: %s with params: page=%d, size=%d, search=%s",
		roleID, params.PageNumber, params.PageSize, params.Search))

	var permissions []entity.Permission
	var totalItems int64

	query := r.db.WithContext(ctx).
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Where("permissions.is_active = ?", true)

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("permissions.name ILIKE ? OR permissions.code ILIKE ? OR permissions.action ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		logger.Error("repository", "GetPermissionsByRoleID", fmt.Sprintf("Failed to count permissions: %v", err))
		return entity.PaginatedPermissions{}, err
	}

	// Apply pagination
	offset := (params.PageNumber - 1) * params.PageSize
	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("permissions.created_at DESC").
		Find(&permissions).Error; err != nil {
		logger.Error("repository", "GetPermissionsByRoleID", fmt.Sprintf("Failed to get permissions: %v", err))
		return entity.PaginatedPermissions{}, err
	}

	result := entity.PaginatedPermissions{
		Items:      permissions,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   params.PageSize,
	}

	logger.Info("repository", "GetPermissionsByRoleID", fmt.Sprintf("Found %d permissions", len(permissions)))
	return result, nil
}

// AssignRolesToUser gán nhiều roles cho một user
func (r *permissionRepository) AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error {
	logger.Info("repository", "AssignRolesToUser", fmt.Sprintf("Assigning %d roles to user: %s", len(roleIDs), userID))

	userRoles := make([]entity.UserRole, len(roleIDs))
	for i, roleID := range roleIDs {
		userRoles[i] = entity.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
	}

	err := r.db.WithContext(ctx).Create(&userRoles).Error
	if err != nil {
		logger.Error("repository", "AssignRolesToUser", fmt.Sprintf("Failed to assign roles: %v", err))
		return err
	}

	logger.Info("repository", "AssignRolesToUser", "Roles assigned successfully")
	return nil
}

// RemoveRolesFromUser xóa roles khỏi user
func (r *permissionRepository) RemoveRolesFromUser(ctx context.Context, userID string, roleIDs []string) error {
	logger.Info("repository", "RemoveRolesFromUser", fmt.Sprintf("Removing %d roles from user: %s", len(roleIDs), userID))

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&entity.UserRole{}).Error

	if err != nil {
		logger.Error("repository", "RemoveRolesFromUser", fmt.Sprintf("Failed to remove roles: %v", err))
		return err
	}

	logger.Info("repository", "RemoveRolesFromUser", "Roles removed successfully")
	return nil
}

// GetRolesByUserID lấy tất cả roles của một user
func (r *permissionRepository) GetRolesByUserID(ctx context.Context, userID string) ([]entity.Role, error) {
	logger.Info("repository", "GetRolesByUserID", fmt.Sprintf("Getting roles for user: %s", userID))

	var roles []entity.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("INNER JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.is_active = ?", true).
		Find(&roles).Error

	if err != nil {
		logger.Error("repository", "GetRolesByUserID", fmt.Sprintf("Failed to get roles: %v", err))
		return nil, err
	}

	logger.Info("repository", "GetRolesByUserID", fmt.Sprintf("Found %d roles", len(roles)))
	return roles, nil
}

// UserHasPermission kiểm tra xem user có permission cụ thể không
func (r *permissionRepository) UserHasPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON role_permissions.role_id = roles.id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permissions.code = ?", permissionCode).
		Where("permissions.is_active = ?", true).
		Where("roles.is_active = ?", true).
		Count(&count).Error

	if err != nil {
		logger.Error("repository", "UserHasPermission", fmt.Sprintf("Failed to check permission: %v", err))
		return false, err
	}

	return count > 0, nil
}

// GetUserPermissions lấy tất cả permissions của một user (từ tất cả roles của user)
func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID string) ([]entity.Permission, error) {
	logger.Info("repository", "GetUserPermissions", fmt.Sprintf("Getting all permissions for user: %s", userID))

	var permissions []entity.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON role_permissions.role_id = roles.id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permissions.is_active = ?", true).
		Where("roles.is_active = ?", true).
		Distinct("permissions.id").
		Find(&permissions).Error

	if err != nil {
		logger.Error("repository", "GetUserPermissions", fmt.Sprintf("Failed to get permissions: %v", err))
		return nil, err
	}

	logger.Info("repository", "GetUserPermissions", fmt.Sprintf("Found %d permissions", len(permissions)))
	return permissions, nil
}

