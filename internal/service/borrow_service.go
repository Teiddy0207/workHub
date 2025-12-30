package service

import (
	"context"
	"errors"
	"time"
	"workHub/constant"
	"workHub/internal/dto"
	"workHub/internal/entity"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BorrowService struct {
	BorrowRepo repository.BorrowRepository
	BookRepo   repository.BookRepository
	UserRepo   repository.AuthRepository
	DB         *gorm.DB
}

type BorrowServiceInterface interface {
	BorrowBook(ctx context.Context, userID string, req dto.BorrowRequest) (dto.BorrowResponse, error)
	ReturnBook(ctx context.Context, borrowID string, userID string, userRole string, req dto.ReturnBorrowRequest) (dto.BorrowResponse, error)
	GetByUserID(ctx context.Context, userID string, targetUserID string, userRole string, status string, page params.QueryParams) (dto.PaginatedBorrowResponse, error)
	GetByBookID(ctx context.Context, bookID string, status string, page params.QueryParams) (dto.PaginatedBorrowResponse, error)
	GetActiveBorrows(ctx context.Context, userID string, userRole string, status string, page params.QueryParams) (dto.PaginatedBorrowResponse, error)
}

func NewBorrowService(borrowRepo repository.BorrowRepository, bookRepo repository.BookRepository, userRepo repository.AuthRepository, db *gorm.DB) BorrowServiceInterface {
	return &BorrowService{
		BorrowRepo: borrowRepo,
		BookRepo:   bookRepo,
		UserRepo:   userRepo,
		DB:         db,
	}
}

func (s *BorrowService) BorrowBook(ctx context.Context, userID string, req dto.BorrowRequest) (dto.BorrowResponse, error) {
	// Validate days to borrow
	if req.DaysToBorrow < 1 || req.DaysToBorrow > 30 {
		return dto.BorrowResponse{}, constant.ErrInvalidDaysToBorrow
	}

	// Kiểm tra book tồn tại
	book, err := s.BookRepo.GetByID(ctx, req.BookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BorrowResponse{}, constant.ErrBookNotFound
		}
		return dto.BorrowResponse{}, err
	}

	// Kiểm tra available copies
	if book.AvailableCopies <= 0 {
		return dto.BorrowResponse{}, constant.ErrNoAvailableCopies
	}

	// Kiểm tra user chưa mượn quá 5 sách cùng lúc
	count, err := s.BorrowRepo.CountActiveByUserID(ctx, userID)
	if err != nil {
		return dto.BorrowResponse{}, err
	}
	if count >= 5 {
		return dto.BorrowResponse{}, constant.ErrUserHasMaxBorrows
	}

	// Kiểm tra user chưa mượn cùng book này
	_, err = s.BorrowRepo.GetActiveByUserAndBook(ctx, userID, req.BookID)
	if err == nil {
		return dto.BorrowResponse{}, constant.ErrBookAlreadyBorrowed
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.BorrowResponse{}, err
	}

	// Bắt đầu transaction
	tx := s.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Tạo borrow record
	now := time.Now()
	dueDate := now.AddDate(0, 0, req.DaysToBorrow)
	borrow := &entity.Borrow{
		ID:         uuid.New().String(),
		UserID:     userID,
		BookID:     req.BookID,
		BorrowDate: now,
		DueDate:    dueDate,
		Status:     "borrowed",
		Notes:      nil,
	}

	// Tạo borrow trong transaction
	borrowRepoTx := repository.NewBorrowRepositoryWithDB(tx)
	bookRepoTx := repository.NewBookRepositoryWithDB(tx)

	err = borrowRepoTx.Create(ctx, borrow)
	if err != nil {
		tx.Rollback()
		return dto.BorrowResponse{}, err
	}

	// Giảm available copies
	book.AvailableCopies--
	err = bookRepoTx.Update(ctx, &book)
	if err != nil {
		tx.Rollback()
		return dto.BorrowResponse{}, err
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		return dto.BorrowResponse{}, err
	}

	// Load book info cho response
	borrow.Book = book

	return mapper.ToBorrowResponse(*borrow, false, true), nil
}

func (s *BorrowService) ReturnBook(ctx context.Context, borrowID string, userID string, userRole string, req dto.ReturnBorrowRequest) (dto.BorrowResponse, error) {
	// Lấy borrow record
	borrow, err := s.BorrowRepo.GetByID(ctx, borrowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BorrowResponse{}, constant.ErrBorrowNotFound
		}
		return dto.BorrowResponse{}, err
	}

	// Kiểm tra quyền: chỉ user mượn sách đó hoặc admin
	if userRole != "admin" && borrow.UserID != userID {
		return dto.BorrowResponse{}, constant.ErrForbidden
	}

	// Kiểm tra đã trả chưa
	if borrow.Status == "returned" {
		return dto.BorrowResponse{}, constant.ErrBorrowAlreadyReturned
	}

	// Bắt đầu transaction
	tx := s.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cập nhật borrow
	now := time.Now()
	borrow.ReturnDate = &now
	borrow.Status = "returned"
	if req.Notes != nil {
		borrow.Notes = req.Notes
	}

	borrowRepoTx := repository.NewBorrowRepositoryWithDB(tx)
	bookRepoTx := repository.NewBookRepositoryWithDB(tx)

	err = borrowRepoTx.Update(ctx, &borrow)
	if err != nil {
		tx.Rollback()
		return dto.BorrowResponse{}, err
	}

	// Tăng available copies
	book, err := bookRepoTx.GetByID(ctx, borrow.BookID)
	if err != nil {
		tx.Rollback()
		return dto.BorrowResponse{}, err
	}

	book.AvailableCopies++
	err = bookRepoTx.Update(ctx, &book)
	if err != nil {
		tx.Rollback()
		return dto.BorrowResponse{}, err
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		return dto.BorrowResponse{}, err
	}

	// Load book info
	borrow.Book = book

	return mapper.ToBorrowResponse(borrow, false, true), nil
}

func (s *BorrowService) GetByUserID(ctx context.Context, userID string, targetUserID string, userRole string, status string, page params.QueryParams) (dto.PaginatedBorrowResponse, error) {
	// Kiểm tra quyền: admin hoặc chính user đó
	if userRole != "admin" && userID != targetUserID {
		return dto.PaginatedBorrowResponse{}, constant.ErrForbidden
	}

	borrows, err := s.BorrowRepo.GetByUserID(ctx, targetUserID, status, page)
	if err != nil {
		return dto.PaginatedBorrowResponse{}, err
	}

	return mapper.ToPaginatedBorrowResponse(borrows, false, true), nil
}

func (s *BorrowService) GetByBookID(ctx context.Context, bookID string, status string, page params.QueryParams) (dto.PaginatedBorrowResponse, error) {
	borrows, err := s.BorrowRepo.GetByBookID(ctx, bookID, status, page)
	if err != nil {
		return dto.PaginatedBorrowResponse{}, err
	}

	return mapper.ToPaginatedBorrowResponse(borrows, true, false), nil
}

func (s *BorrowService) GetActiveBorrows(ctx context.Context, userID string, userRole string, status string, page params.QueryParams) (dto.PaginatedBorrowResponse, error) {
	// User thường chỉ xem của mình, admin xem tất cả
	queryUserID := ""
	if userRole != "admin" {
		queryUserID = userID
	}

	borrows, err := s.BorrowRepo.GetActiveBorrows(ctx, queryUserID, status, page)
	if err != nil {
		return dto.PaginatedBorrowResponse{}, err
	}

	return mapper.ToPaginatedBorrowResponse(borrows, true, true), nil
}
