package repository

import (
	"context"
	"fmt"
	"workHub/internal/entity"
	"workHub/logger"
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
	logger.Info("repository", "CreateRole", fmt.Sprintf("Creating role: %s (%s)", role.Name, role.Code))

	err := r.db.WithContext(ctx).Create(role).Error
	if err != nil {
		logger.Error("repository", "CreateRole", fmt.Sprintf("Failed to create role: %v", err))
		return err
	}

	logger.Info("repository", "CreateRole", fmt.Sprintf("Role created successfully: %s", role.ID))
	return nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id string) (entity.Role, error) {
	logger.Info("repository", "GetRoleByID", fmt.Sprintf("Getting role by ID: %s", id))

	var role entity.Role
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&role).Error

	if err != nil {
		logger.Error("repository", "GetRoleByID", fmt.Sprintf("Role not found: %s, error: %v", id, err))
		return entity.Role{}, err
	}

	logger.Info("repository", "GetRoleByID", fmt.Sprintf("Role found: %s", role.Name))
	return role, nil
}

func (r *roleRepository) GetRoleByCode(ctx context.Context, code string) (entity.Role, error) {
	logger.Info("repository", "GetRoleByCode", fmt.Sprintf("Getting role by code: %s", code))

	var role entity.Role
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&role).Error

	if err != nil {
		logger.Error("repository", "GetRoleByCode", fmt.Sprintf("Role not found: %s, error: %v", code, err))
		return entity.Role{}, err
	}

	logger.Info("repository", "GetRoleByCode", fmt.Sprintf("Role found: %s", role.Name))
	return role, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, id string, role *entity.Role) error {
	logger.Info("repository", "UpdateRole", fmt.Sprintf("Updating role: %s", id))

	result := r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("id = ?", id).
		Updates(role)

	if result.Error != nil {
		logger.Error("repository", "UpdateRole", fmt.Sprintf("Failed to update role: %s, error: %v", id, result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Error("repository", "UpdateRole", fmt.Sprintf("Role not found for update: %s", id))
		return gorm.ErrRecordNotFound
	}

	logger.Info("repository", "UpdateRole", fmt.Sprintf("Role updated successfully: %s", id))
	return nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, id string) error {
	logger.Info("repository", "DeleteRole", fmt.Sprintf("Hard deleting role: %s", id))

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.Role{})

	if result.Error != nil {
		logger.Error("repository", "DeleteRole", fmt.Sprintf("Failed to delete role: %s, error: %v", id, result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Error("repository", "DeleteRole", fmt.Sprintf("Role not found for deletion: %s", id))
		return gorm.ErrRecordNotFound
	}

	logger.Info("repository", "DeleteRole", fmt.Sprintf("Role deleted successfully: %s", id))
	return nil
}

func (r *roleRepository) ListRoles(ctx context.Context, params params.QueryParams) (entity.PaginatedRoles, error) {
	logger.Info("repository", "ListRoles", fmt.Sprintf("Listing roles with params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

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
		logger.Error("repository", "ListRoles", fmt.Sprintf("Failed to count roles: %v", err))
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
		logger.Error("repository", "ListRoles", fmt.Sprintf("Failed to fetch roles: %v", err))
		return entity.PaginatedRoles{}, err
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10 // Giá trị mặc định
	}

	logger.Info("repository", "ListRoles", fmt.Sprintf("Found %d roles (total: %d)", len(roles), totalItems))

	// Trả kết quả
	return entity.PaginatedRoles{
		Items:      roles,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   pageSize,
	}, nil
}
