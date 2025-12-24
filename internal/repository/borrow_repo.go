package repository

import (
	"context"
	"workHub/internal/entity"
	"workHub/pkg/params"

	"gorm.io/gorm"
)

type borrowRepository struct {
	db *gorm.DB
}

// NewBorrowRepositoryWithDB tạo repository với DB instance (để dùng trong transaction)
func NewBorrowRepositoryWithDB(db *gorm.DB) BorrowRepository {
	return &borrowRepository{db: db}
}

type BorrowRepository interface {
	Create(ctx context.Context, borrow *entity.Borrow) error
	GetByID(ctx context.Context, id string) (entity.Borrow, error)
	GetByUserID(ctx context.Context, userID string, status string, page params.QueryParams) (entity.PaginatedBorrows, error)
	GetByBookID(ctx context.Context, bookID string, status string, page params.QueryParams) (entity.PaginatedBorrows, error)
	GetActiveByUserID(ctx context.Context, userID string) ([]entity.Borrow, error)
	GetActiveByUserAndBook(ctx context.Context, userID string, bookID string) (entity.Borrow, error)
	Update(ctx context.Context, borrow *entity.Borrow) error
	GetActiveBorrows(ctx context.Context, userID string, status string, page params.QueryParams) (entity.PaginatedBorrows, error)
	CountActiveByUserID(ctx context.Context, userID string) (int64, error)
}

func NewBorrowRepository(db *gorm.DB) BorrowRepository {
	return &borrowRepository{db: db}
}

func (r *borrowRepository) Create(ctx context.Context, borrow *entity.Borrow) error {
	return r.db.WithContext(ctx).Create(borrow).Error
}

func (r *borrowRepository) GetByID(ctx context.Context, id string) (entity.Borrow, error) {
	var borrow entity.Borrow
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Book").
		Where("id = ?", id).
		First(&borrow).Error
	return borrow, err
}

func (r *borrowRepository) GetByUserID(ctx context.Context, userID string, status string, page params.QueryParams) (entity.PaginatedBorrows, error) {
	var borrows []entity.Borrow
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Borrow{}).
		Preload("Book").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedBorrows{}, err
	}

	offset := (page.PageNumber - 1) * page.PageSize
	pageSize := page.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("borrow_date DESC").
		Find(&borrows).Error; err != nil {
		return entity.PaginatedBorrows{}, err
	}

	return entity.PaginatedBorrows{
		Items:      borrows,
		TotalItems: int(totalItems),
		PageNumber: page.PageNumber,
		PageSize:   pageSize,
	}, nil
}

func (r *borrowRepository) GetByBookID(ctx context.Context, bookID string, status string, page params.QueryParams) (entity.PaginatedBorrows, error) {
	var borrows []entity.Borrow
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Borrow{}).
		Preload("User").
		Where("book_id = ?", bookID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedBorrows{}, err
	}

	offset := (page.PageNumber - 1) * page.PageSize
	pageSize := page.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("borrow_date DESC").
		Find(&borrows).Error; err != nil {
		return entity.PaginatedBorrows{}, err
	}

	return entity.PaginatedBorrows{
		Items:      borrows,
		TotalItems: int(totalItems),
		PageNumber: page.PageNumber,
		PageSize:   pageSize,
	}, nil
}

func (r *borrowRepository) GetActiveByUserID(ctx context.Context, userID string) ([]entity.Borrow, error) {
	var borrows []entity.Borrow
	err := r.db.WithContext(ctx).
		Preload("Book").
		Where("user_id = ? AND status != ?", userID, "returned").
		Find(&borrows).Error
	return borrows, err
}

func (r *borrowRepository) GetActiveByUserAndBook(ctx context.Context, userID string, bookID string) (entity.Borrow, error) {
	var borrow entity.Borrow
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND book_id = ? AND status != ?", userID, bookID, "returned").
		First(&borrow).Error
	return borrow, err
}

func (r *borrowRepository) Update(ctx context.Context, borrow *entity.Borrow) error {
	return r.db.WithContext(ctx).Save(borrow).Error
}

func (r *borrowRepository) GetActiveBorrows(ctx context.Context, userID string, status string, page params.QueryParams) (entity.PaginatedBorrows, error) {
	var borrows []entity.Borrow
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Borrow{}).
		Preload("User").
		Preload("Book").
		Where("status != ?", "returned")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		statuses := []string{status}
		if status == "overdue" {
			query = query.Where("due_date < CURRENT_DATE")
		} else {
			query = query.Where("status IN ?", statuses)
		}
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedBorrows{}, err
	}

	offset := (page.PageNumber - 1) * page.PageSize
	pageSize := page.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("borrow_date DESC").
		Find(&borrows).Error; err != nil {
		return entity.PaginatedBorrows{}, err
	}

	return entity.PaginatedBorrows{
		Items:      borrows,
		TotalItems: int(totalItems),
		PageNumber: page.PageNumber,
		PageSize:   pageSize,
	}, nil
}

func (r *borrowRepository) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Borrow{}).
		Where("user_id = ? AND status != ?", userID, "returned").
		Count(&count).Error
	return count, err
}

