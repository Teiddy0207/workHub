package repository

import (
	"context"
	"fmt"
	"workHub/internal/entity"
	"workHub/pkg/params"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

type RoleRepository interface {
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRoleByID(ctx context.Context, id string) (entity.Role, error)
	GetRoleByCode(ctx context.Context, code string) (entity.Role, error)
	UpdateRole(ctx context.Context, id string, role *entity.Role) error
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context, params params.QueryParams) (entity.PaginatedRoles, error)
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) CreateRole(ctx context.Context, role *entity.Role) error {
	fmt.Printf("🔍 Creating role: %s (%s)\n", role.Name, role.Code)
	
	err := r.db.WithContext(ctx).Create(role).Error
	if err != nil {
		fmt.Printf("❌ Failed to create role: %v\n", err)
		return err
	}
	
	fmt.Printf("✅ Role created successfully: %s\n", role.ID)
	return nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id string) (entity.Role, error) {
	fmt.Printf("🔍 Getting role by ID: %s\n", id)
	
	var role entity.Role
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&role).Error
	
	if err != nil {
		fmt.Printf("❌ Role not found: %v\n", err)
		return entity.Role{}, err
	}
	
	fmt.Printf("✅ Role found: %s\n", role.Name)
	return role, nil
}

func (r *roleRepository) GetRoleByCode(ctx context.Context, code string) (entity.Role, error) {
	fmt.Printf("🔍 Getting role by code: %s\n", code)
	
	var role entity.Role
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&role).Error
	
	if err != nil {
		fmt.Printf("❌ Role not found: %v\n", err)
		return entity.Role{}, err
	}
	
	fmt.Printf("✅ Role found: %s\n", role.Name)
	return role, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, id string, role *entity.Role) error {
	fmt.Printf("🔍 Updating role: %s\n", id)
	
	result := r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("id = ?", id).
		Updates(role)
	
	if result.Error != nil {
		fmt.Printf("❌ Failed to update role: %v\n", result.Error)
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		fmt.Printf("❌ Role not found for update: %s\n", id)
		return gorm.ErrRecordNotFound
	}
	
	fmt.Printf("✅ Role updated successfully: %s\n", id)
	return nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, id string) error {
	fmt.Printf("🔍 Hard deleting role: %s\n", id)
	
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.Role{})
	
	if result.Error != nil {
		fmt.Printf("❌ Failed to delete role: %v\n", result.Error)
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		fmt.Printf("❌ Role not found for deletion: %s\n", id)
		return gorm.ErrRecordNotFound
	}
	
	fmt.Printf("✅ Role deleted successfully: %s\n", id)
	return nil
}

func (r *roleRepository) ListRoles(ctx context.Context, params params.QueryParams) (entity.PaginatedRoles, error) {
	fmt.Printf("🔍 Listing roles with params: page=%d, size=%d, search=%s\n", 
		params.PageNumber, params.PageSize, params.Search)
	
	var roles []entity.Role
	var totalItems int64

	// Tạo query cơ bản
	query := r.db.WithContext(ctx).Model(&entity.Role{})

	// Nếu có tìm kiếm theo name hoặc code
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", 
			searchTerm, searchTerm, searchTerm)
	}

	// Đếm tổng số role trước khi phân trang
	if err := query.Count(&totalItems).Error; err != nil {
		fmt.Printf("❌ Failed to count roles: %v\n", err)
		return entity.PaginatedRoles{}, err
	}

	// Tính offset
	offset := (params.PageNumber - 1) * params.PageSize

	// Lấy dữ liệu phân trang
	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&roles).Error; err != nil {
		fmt.Printf("❌ Failed to fetch roles: %v\n", err)
		return entity.PaginatedRoles{}, err
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10 // Giá trị mặc định
	}

	fmt.Printf("✅ Found %d roles (total: %d)\n", len(roles), totalItems)

	// Trả kết quả
	return entity.PaginatedRoles{
		Items:      roles,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   pageSize,
	}, nil
}
