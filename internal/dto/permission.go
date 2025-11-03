package dto

import "time"

// PermissionRequest - DTO cho request tạo permission
type PermissionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"required,min=2,max=100"`
	Action      string `json:"action" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=500"`
}

// PermissionResponse - DTO cho response permission
type PermissionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionUpdateRequest - DTO cho request cập nhật permission
type PermissionUpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Code        *string `json:"code" binding:"omitempty,min=2,max=100"`
	Action      *string `json:"action" binding:"omitempty,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

// AssignPermissionRequest - DTO cho request gán permission cho role
type AssignPermissionRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1"`
}

// AssignRoleRequest - DTO cho request gán role cho user
type AssignRoleRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required,min=1"`
}

// RoleWithPermissionsResponse - DTO cho response role kèm danh sách permissions có phân trang
type RoleWithPermissionsResponse struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Code        string                       `json:"code"`
	Description string                       `json:"description"`
	IsActive    bool                         `json:"is_active"`
	Permissions PaginatedPermissionResponse  `json:"permissions"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
}

// PaginatedPermissionResponse - DTO cho response danh sách permission có phân trang
type PaginatedPermissionResponse = Pagination[PermissionResponse]

