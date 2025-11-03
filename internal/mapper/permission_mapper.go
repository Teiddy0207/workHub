package mapper

import (
	"workHub/internal/dto"
	"workHub/internal/entity"
	"github.com/google/uuid"
)

// ToPermissionResponse chuyển đổi entity.Permission thành dto.PermissionResponse
func ToPermissionResponse(permission entity.Permission) dto.PermissionResponse {
	return dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Code:        permission.Code,
		Action:      permission.Action,
		Description: permission.Description,
		IsActive:    permission.IsActive,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}

// ToPermissionResponseList chuyển đổi danh sách entity.Permission thành danh sách dto.PermissionResponse
func ToPermissionResponseList(permissions []entity.Permission) []dto.PermissionResponse {
	responses := make([]dto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = ToPermissionResponse(permission)
	}
	return responses
}

// ToPaginatedPermissionResponse chuyển đổi entity.PaginatedPermissions thành dto.PaginatedPermissionResponse
func ToPaginatedPermissionResponse(paginatedPermissions entity.PaginatedPermissions) dto.PaginatedPermissionResponse {
	return dto.PaginatedPermissionResponse{
		Items:       ToPermissionResponseList(paginatedPermissions.Items),
		TotalItems:  paginatedPermissions.TotalItems,
		TotalPages:  (paginatedPermissions.TotalItems + paginatedPermissions.PageSize - 1) / paginatedPermissions.PageSize,
		CurrentPage: paginatedPermissions.PageNumber,
		PageSize:    paginatedPermissions.PageSize,
	}
}

// ToPermissionEntity chuyển đổi dto.PermissionRequest thành entity.Permission
func ToPermissionEntity(req dto.PermissionRequest) entity.Permission {
	return entity.Permission{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Code:        req.Code,
		Action:      req.Action,
		Description: req.Description,
		IsActive:    true, // Mặc định là true
	}
}

// ToPermissionUpdateEntity chuyển đổi dto.PermissionUpdateRequest thành entity.Permission (chỉ các field cần update)
func ToPermissionUpdateEntity(req dto.PermissionUpdateRequest) *entity.Permission {
	updateData := &entity.Permission{}
	
	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Code != nil {
		updateData.Code = *req.Code
	}
	if req.Action != nil {
		updateData.Action = *req.Action
	}
	if req.Description != nil {
		updateData.Description = *req.Description
	}
	
	return updateData
}

// ToRoleWithPermissionsResponse chuyển đổi entity.Role và entity.PaginatedPermissions thành dto.RoleWithPermissionsResponse
func ToRoleWithPermissionsResponse(role entity.Role, permissions entity.PaginatedPermissions) dto.RoleWithPermissionsResponse {
	return dto.RoleWithPermissionsResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		IsActive:    role.IsActive,
		Permissions: ToPaginatedPermissionResponse(permissions),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

