package dto

import "time"

// RoleRequest - DTO cho request tạo/cập nhật role
type RoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=500"`
	IsActive    *bool  `json:"is_active"`
}

// RoleResponse - DTO cho response role
type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoleUpdateRequest - DTO cho request cập nhật role
type RoleUpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Code        *string `json:"code" binding:"omitempty,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	IsActive    *bool   `json:"is_active"`
}

// PaginatedRoleResponse - DTO cho response danh sách role có phân trang
type PaginatedRoleResponse = Pagination[RoleResponse]
