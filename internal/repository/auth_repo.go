package repository

import (
	"context"
	"workHub/internal/entity"

	// "workHub/internal/dto"
	// "workHub/internal/mapper"
	// "workHub/internal/repository"
	// "workHub/internal/entity"

	"gorm.io/gorm"
	// "workHub/internal/service"
	"workHub/pkg/params"
)

type authRepository struct {
	db *gorm.DB
}

type AuthRepository interface {
	ListUsers(ctx context.Context, params params.QueryParams) (entity.PaginatedUsers, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) ListUsers(ctx context.Context, params params.QueryParams) (entity.PaginatedUsers, error) {
	var users []entity.User
	var totalItems int64

	// Tạo query cơ bản
	query := r.db.WithContext(ctx).Model(&entity.User{})

	// Nếu có tìm kiếm theo username hoặc email
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Đếm tổng số user trước khi phân trang
	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedUsers{}, err
	}

	// Tính offset
	offset := (params.PageNumber - 1) * params.PageSize

	// Lấy dữ liệu phân trang
	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return entity.PaginatedUsers{}, err
	}
	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedUsers{}, err
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10 // Giá trị mặc định
	}
	// Trả kết quả
	return entity.PaginatedUsers{
		Items:      users,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   pageSize,
	}, nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	
	if err != nil {
		return entity.User{}, err
	}
	
	return user, nil
}
