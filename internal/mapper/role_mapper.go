package mapper

import (
	"workHub/internal/dto"
	"workHub/internal/entity"
	"github.com/google/uuid"
)

// ToRoleResponse chuyển đổi entity.Role thành dto.RoleResponse
func ToRoleResponse(role entity.Role) dto.RoleResponse {
	return dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// ToRoleResponseList chuyển đổi danh sách entity.Role thành danh sách dto.RoleResponse
func ToRoleResponseList(roles []entity.Role) []dto.RoleResponse {
	responses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToRoleResponse(role)
	}
	return responses
}

// ToPaginatedRoleResponse chuyển đổi entity.PaginatedRoles thành dto.PaginatedRoleResponse
func ToPaginatedRoleResponse(paginatedRoles entity.PaginatedRoles) dto.PaginatedRoleResponse {
	return dto.PaginatedRoleResponse{
		Items:       ToRoleResponseList(paginatedRoles.Items),
		TotalItems:  paginatedRoles.TotalItems,
		TotalPages:  (paginatedRoles.TotalItems + paginatedRoles.PageSize - 1) / paginatedRoles.PageSize,
		CurrentPage: paginatedRoles.PageNumber,
		PageSize:    paginatedRoles.PageSize,
	}
}

// ToRoleEntity chuyển đổi dto.RoleRequest thành entity.Role
func ToRoleEntity(req dto.RoleRequest) entity.Role {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return entity.Role{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    isActive,
	}
}

// ToRoleUpdateEntity chuyển đổi dto.RoleUpdateRequest thành entity.Role (chỉ các field cần update)
func ToRoleUpdateEntity(req dto.RoleUpdateRequest) *entity.Role {
	updateData := &entity.Role{}
	
	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Code != nil {
		updateData.Code = *req.Code
	}
	if req.Description != nil {
		updateData.Description = *req.Description
	}
	if req.IsActive != nil {
		updateData.IsActive = *req.IsActive
	}
	
	return updateData
}
