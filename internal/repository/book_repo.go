package repository

import (
	"context"
	"workHub/internal/entity"
	"workHub/pkg/params"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepositoryWithDB tạo repository với DB instance (để dùng trong transaction)
func NewBookRepositoryWithDB(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	GetByID(ctx context.Context, id string) (entity.Book, error)
	GetByISBN(ctx context.Context, isbn string) (entity.Book, error)
	List(ctx context.Context, params params.QueryParams, category string, search string) (entity.PaginatedBooks, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, searchReq entity.BookSearchRequest, page params.QueryParams) (entity.PaginatedBooks, error)
	HasActiveBorrows(ctx context.Context, bookID string) (bool, error)
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) GetByID(ctx context.Context, id string) (entity.Book, error) {
	var book entity.Book
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&book).Error
	return book, err
}

func (r *bookRepository) GetByISBN(ctx context.Context, isbn string) (entity.Book, error) {
	var book entity.Book
	err := r.db.WithContext(ctx).Where("isbn = ?", isbn).First(&book).Error
	return book, err
}

func (r *bookRepository) List(ctx context.Context, params params.QueryParams, category string, search string) (entity.PaginatedBooks, error) {
	var books []entity.Book
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Book{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedBooks{}, err
	}

	offset := (params.PageNumber - 1) * params.PageSize
	sortBy := "created_at"
	sortOrder := "DESC"


	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order(sortBy + " " + sortOrder).
		Find(&books).Error; err != nil {
		return entity.PaginatedBooks{}, err
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	return entity.PaginatedBooks{
		Items:      books,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   pageSize,
	}, nil
}

func (r *bookRepository) Update(ctx context.Context, book *entity.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Book{}, "id = ?", id).Error
}

func (r *bookRepository) Search(ctx context.Context, searchReq entity.BookSearchRequest, page params.QueryParams) (entity.PaginatedBooks, error) {
	var books []entity.Book
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Book{})

	if searchReq.Q != "" {
		searchTerm := "%" + searchReq.Q + "%"
		query = query.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if searchReq.Category != "" {
		query = query.Where("category = ?", searchReq.Category)
	}

	if searchReq.Author != "" {
		query = query.Where("author ILIKE ?", "%"+searchReq.Author+"%")
	}

	if searchReq.AvailableOnly {
		query = query.Where("available_copies > 0")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return entity.PaginatedBooks{}, err
	}

	offset := (page.PageNumber - 1) * page.PageSize
	pageSize := page.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&books).Error; err != nil {
		return entity.PaginatedBooks{}, err
	}

	return entity.PaginatedBooks{
		Items:      books,
		TotalItems: int(totalItems),
		PageNumber: page.PageNumber,
		PageSize:   pageSize,
	}, nil
}

func (r *bookRepository) HasActiveBorrows(ctx context.Context, bookID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Borrow{}).
		Where("book_id = ? AND status != ?", bookID, "returned").
		Count(&count).Error
	return count > 0, err
}
