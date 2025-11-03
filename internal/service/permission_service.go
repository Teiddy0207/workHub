package service

import (
	"context"
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/constant"
	"workHub/logger"
)

type PermissionService struct {
	permissionRepo repository.PermissionRepository
	roleRepo       repository.RoleRepository
}

type PermissionServiceInterface interface {
	CreatePermission(ctx context.Context, req dto.PermissionRequest) (dto.PermissionResponse, error)
	GetPermissionByID(ctx context.Context, id string) (dto.PermissionResponse, error)
	UpdatePermission(ctx context.Context, id string, req dto.PermissionUpdateRequest) (dto.PermissionResponse, error)
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context, params params.QueryParams) (dto.PaginatedPermissionResponse, error)
	
	// Role-Permission management
	AssignPermissionsToRole(ctx context.Context, roleID string, req dto.AssignPermissionRequest) error
	RemovePermissionsFromRole(ctx context.Context, roleID string, req dto.AssignPermissionRequest) error
	GetRoleWithPermissions(ctx context.Context, roleID string, params params.QueryParams) (dto.RoleWithPermissionsResponse, error)
	
	// User-Role management
	AssignRolesToUser(ctx context.Context, userID string, req dto.AssignRoleRequest) error
	RemoveRolesFromUser(ctx context.Context, userID string, req dto.AssignRoleRequest) error
	GetUserPermissions(ctx context.Context, userID string) ([]dto.PermissionResponse, error)
	UserHasPermission(ctx context.Context, userID string, permissionCode string) (bool, error)
}

func NewPermissionService(permissionRepo repository.PermissionRepository, roleRepo repository.RoleRepository) PermissionServiceInterface {
	return &PermissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
	}
}

func (s *PermissionService) CreatePermission(ctx context.Context, req dto.PermissionRequest) (dto.PermissionResponse, error) {
	logger.Info("service", "CreatePermission", fmt.Sprintf("Creating new permission: %s (%s)", req.Name, req.Code))

	// Kiểm tra code đã tồn tại chưa
	_, err := s.permissionRepo.GetPermissionByCode(ctx, req.Code)
	if err == nil {
		logger.Warn("service", "CreatePermission", fmt.Sprintf("Permission code already exists: %s", req.Code))
		return dto.PermissionResponse{}, constant.ErrTakenCredential
	}

	permission := mapper.ToPermissionEntity(req)

	// Tạo permission trong database
	err = s.permissionRepo.CreatePermission(ctx, &permission)
	if err != nil {
		logger.Error("service", "CreatePermission", fmt.Sprintf("Failed to create permission: %v", err))
		return dto.PermissionResponse{}, err
	}

	createdPermission, err := s.permissionRepo.GetPermissionByID(ctx, permission.ID)
	if err != nil {
		logger.Error("service", "CreatePermission", fmt.Sprintf("Failed to get created permission: %v", err))
		return dto.PermissionResponse{}, err
	}

	logger.Info("service", "CreatePermission", fmt.Sprintf("Permission created successfully: %s", createdPermission.Name))
	return mapper.ToPermissionResponse(createdPermission), nil
}

func (s *PermissionService) GetPermissionByID(ctx context.Context, id string) (dto.PermissionResponse, error) {
	logger.Info("service", "GetPermissionByID", fmt.Sprintf("Getting permission by ID: %s", id))

	permission, err := s.permissionRepo.GetPermissionByID(ctx, id)
	if err != nil {
		logger.Error("service", "GetPermissionByID", fmt.Sprintf("Permission not found: %s, error: %v", id, err))
		return dto.PermissionResponse{}, constant.ErrNotFound
	}

	logger.Info("service", "GetPermissionByID", fmt.Sprintf("Permission found: %s", permission.Name))
	return mapper.ToPermissionResponse(permission), nil
}

func (s *PermissionService) UpdatePermission(ctx context.Context, id string, req dto.PermissionUpdateRequest) (dto.PermissionResponse, error) {
	logger.Info("service", "UpdatePermission", fmt.Sprintf("Updating permission: %s", id))

	existingPermission, err := s.permissionRepo.GetPermissionByID(ctx, id)
	if err != nil {
		logger.Error("service", "UpdatePermission", fmt.Sprintf("Permission not found: %s, error: %v", id, err))
		return dto.PermissionResponse{}, constant.ErrNotFound
	}

	// Kiểm tra code nếu có thay đổi
	if req.Code != nil && *req.Code != existingPermission.Code {
		if _, err := s.permissionRepo.GetPermissionByCode(ctx, *req.Code); err == nil {
			logger.Warn("service", "UpdatePermission", fmt.Sprintf("Permission code already exists: %s", *req.Code))
			return dto.PermissionResponse{}, constant.ErrTakenCredential
		}
	}

	updateData := mapper.ToPermissionUpdateEntity(req)
	if err := s.permissionRepo.UpdatePermission(ctx, id, updateData); err != nil {
		logger.Error("service", "UpdatePermission", fmt.Sprintf("Failed to update permission: %s, error: %v", id, err))
		return dto.PermissionResponse{}, err
	}

	updatedPermission, err := s.permissionRepo.GetPermissionByID(ctx, id)
	if err != nil {
		logger.Error("service", "UpdatePermission", fmt.Sprintf("Failed to get updated permission: %s, error: %v", id, err))
		return dto.PermissionResponse{}, err
	}

	logger.Info("service", "UpdatePermission", fmt.Sprintf("Permission updated successfully: %s", updatedPermission.Name))
	return mapper.ToPermissionResponse(updatedPermission), nil
}

func (s *PermissionService) DeletePermission(ctx context.Context, id string) error {
	logger.Info("service", "DeletePermission", fmt.Sprintf("Deleting permission: %s", id))

	_, err := s.permissionRepo.GetPermissionByID(ctx, id)
	if err != nil {
		logger.Error("service", "DeletePermission", fmt.Sprintf("Permission not found: %s, error: %v", id, err))
		return constant.ErrNotFound
	}

	err = s.permissionRepo.DeletePermission(ctx, id)
	if err != nil {
		logger.Error("service", "DeletePermission", fmt.Sprintf("Failed to delete permission: %s, error: %v", id, err))
		return err
	}

	logger.Info("service", "DeletePermission", fmt.Sprintf("Permission deleted successfully: %s", id))
	return nil
}

func (s *PermissionService) ListPermissions(ctx context.Context, params params.QueryParams) (dto.PaginatedPermissionResponse, error) {
	logger.Info("service", "ListPermissions", fmt.Sprintf("Listing permissions with params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	permissions, err := s.permissionRepo.ListPermissions(ctx, params)
	if err != nil {
		logger.Error("service", "ListPermissions", fmt.Sprintf("Failed to list permissions: %v", err))
		return dto.PaginatedPermissionResponse{}, err
	}

	logger.Info("service", "ListPermissions", fmt.Sprintf("Listed %d permissions successfully", len(permissions.Items)))
	return mapper.ToPaginatedPermissionResponse(permissions), nil
}

func (s *PermissionService) AssignPermissionsToRole(ctx context.Context, roleID string, req dto.AssignPermissionRequest) error {
	logger.Info("service", "AssignPermissionsToRole", fmt.Sprintf("Assigning permissions to role: %s", roleID))

	// Kiểm tra role tồn tại
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		logger.Error("service", "AssignPermissionsToRole", fmt.Sprintf("Role not found: %s", roleID))
		return constant.ErrNotFound
	}

	// Kiểm tra tất cả permissions tồn tại
	for _, permissionID := range req.PermissionIDs {
		_, err := s.permissionRepo.GetPermissionByID(ctx, permissionID)
		if err != nil {
			logger.Error("service", "AssignPermissionsToRole", fmt.Sprintf("Permission not found: %s", permissionID))
			return constant.ErrNotFound
		}
	}

	err = s.permissionRepo.AssignPermissionsToRole(ctx, roleID, req.PermissionIDs)
	if err != nil {
		logger.Error("service", "AssignPermissionsToRole", fmt.Sprintf("Failed to assign permissions: %v", err))
		return err
	}

	logger.Info("service", "AssignPermissionsToRole", "Permissions assigned successfully")
	return nil
}

func (s *PermissionService) RemovePermissionsFromRole(ctx context.Context, roleID string, req dto.AssignPermissionRequest) error {
	logger.Info("service", "RemovePermissionsFromRole", fmt.Sprintf("Removing permissions from role: %s", roleID))

	// Kiểm tra role tồn tại
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		logger.Error("service", "RemovePermissionsFromRole", fmt.Sprintf("Role not found: %s", roleID))
		return constant.ErrNotFound
	}

	err = s.permissionRepo.RemovePermissionsFromRole(ctx, roleID, req.PermissionIDs)
	if err != nil {
		logger.Error("service", "RemovePermissionsFromRole", fmt.Sprintf("Failed to remove permissions: %v", err))
		return err
	}

	logger.Info("service", "RemovePermissionsFromRole", "Permissions removed successfully")
	return nil
}

func (s *PermissionService) GetRoleWithPermissions(ctx context.Context, roleID string, params params.QueryParams) (dto.RoleWithPermissionsResponse, error) {
	logger.Info("service", "GetRoleWithPermissions", fmt.Sprintf("Getting role with permissions: %s", roleID))

	// Lấy role
	role, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		logger.Error("service", "GetRoleWithPermissions", fmt.Sprintf("Role not found: %s", roleID))
		return dto.RoleWithPermissionsResponse{}, constant.ErrNotFound
	}

	// Lấy permissions của role có phân trang
	permissions, err := s.permissionRepo.GetPermissionsByRoleID(ctx, roleID, params)
	if err != nil {
		logger.Error("service", "GetRoleWithPermissions", fmt.Sprintf("Failed to get permissions: %v", err))
		return dto.RoleWithPermissionsResponse{}, err
	}

	response := mapper.ToRoleWithPermissionsResponse(role, permissions)

	logger.Info("service", "GetRoleWithPermissions", fmt.Sprintf("Role with %d permissions retrieved", len(permissions.Items)))
	return response, nil
}

func (s *PermissionService) AssignRolesToUser(ctx context.Context, userID string, req dto.AssignRoleRequest) error {
	logger.Info("service", "AssignRolesToUser", fmt.Sprintf("Assigning roles to user: %s", userID))

	// Kiểm tra tất cả roles tồn tại
	for _, roleID := range req.RoleIDs {
		_, err := s.roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			logger.Error("service", "AssignRolesToUser", fmt.Sprintf("Role not found: %s", roleID))
			return constant.ErrNotFound
		}
	}

	err := s.permissionRepo.AssignRolesToUser(ctx, userID, req.RoleIDs)
	if err != nil {
		logger.Error("service", "AssignRolesToUser", fmt.Sprintf("Failed to assign roles: %v", err))
		return err
	}

	logger.Info("service", "AssignRolesToUser", "Roles assigned successfully")
	return nil
}

func (s *PermissionService) RemoveRolesFromUser(ctx context.Context, userID string, req dto.AssignRoleRequest) error {
	logger.Info("service", "RemoveRolesFromUser", fmt.Sprintf("Removing roles from user: %s", userID))

	err := s.permissionRepo.RemoveRolesFromUser(ctx, userID, req.RoleIDs)
	if err != nil {
		logger.Error("service", "RemoveRolesFromUser", fmt.Sprintf("Failed to remove roles: %v", err))
		return err
	}

	logger.Info("service", "RemoveRolesFromUser", "Roles removed successfully")
	return nil
}

func (s *PermissionService) GetUserPermissions(ctx context.Context, userID string) ([]dto.PermissionResponse, error) {
	logger.Info("service", "GetUserPermissions", fmt.Sprintf("Getting permissions for user: %s", userID))

	permissions, err := s.permissionRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("service", "GetUserPermissions", fmt.Sprintf("Failed to get permissions: %v", err))
		return nil, err
	}

	logger.Info("service", "GetUserPermissions", fmt.Sprintf("Found %d permissions", len(permissions)))
	return mapper.ToPermissionResponseList(permissions), nil
}

func (s *PermissionService) UserHasPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	return s.permissionRepo.UserHasPermission(ctx, userID, permissionCode)
}

