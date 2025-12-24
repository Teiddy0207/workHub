package repository

import (
	"context"
	"workHub/internal/entity"
	"workHub/pkg/params"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

type AuthRepository interface {
	ListUsers(ctx context.Context, params params.QueryParams) (entity.PaginatedUsers, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
	HasActiveBorrows(ctx context.Context, userID string) (bool, error)
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) ListUsers(ctx context.Context, params params.QueryParams) (entity.PaginatedUsers, error) {
	var users []entity.User
	var totalItems int64

	// Tạo query cơ bản
	query := r.db.WithContext(ctx).Model(&entity.User{})

	// Nếu có tìm kiếm theo email hoặc full_name
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("email ILIKE ? OR full_name ILIKE ?", searchTerm, searchTerm)
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

func (r *authRepository) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error
	return user, err
}

func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *authRepository) HasActiveBorrows(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Borrow{}).
		Where("user_id = ? AND status != ?", userID, "returned").
		Count(&count).Error
	return count > 0, err
}
