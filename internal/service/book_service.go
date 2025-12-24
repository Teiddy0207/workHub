package service

import (
	"context"
	"errors"
	"workHub/internal/dto"
	"workHub/internal/entity"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/constant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookService struct {
	BookRepo repository.BookRepository
}

type BookServiceInterface interface {
	Create(ctx context.Context, req dto.BookRequest) (dto.BookResponse, error)
	GetByID(ctx context.Context, id string) (dto.BookResponse, error)
	List(ctx context.Context, page params.QueryParams, category string, search string) (dto.PaginatedBookResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateBookRequest) (dto.BookResponse, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, searchReq dto.BookSearchRequest, page params.QueryParams) (dto.PaginatedBookResponse, error)
}

func NewBookService(bookRepo repository.BookRepository) BookServiceInterface {
	return &BookService{
		BookRepo: bookRepo,
	}
}

func (s *BookService) Create(ctx context.Context, req dto.BookRequest) (dto.BookResponse, error) {
	// Kiểm tra ISBN đã tồn tại chưa
	_, err := s.BookRepo.GetByISBN(ctx, req.ISBN)
	if err == nil {
		return dto.BookResponse{}, constant.ErrISBNAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.BookResponse{}, err
	}

	book := entity.Book{
		ID:              uuid.New().String(),
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		Category:        req.Category,
		Description:     req.Description,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies,
		PublishedYear:   req.PublishedYear,
		Publisher:       req.Publisher,
	}

	err = s.BookRepo.Create(ctx, &book)
	if err != nil {
		return dto.BookResponse{}, err
	}

	return mapper.ToBookResponse(book), nil
}

func (s *BookService) GetByID(ctx context.Context, id string) (dto.BookResponse, error) {
	book, err := s.BookRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BookResponse{}, constant.ErrBookNotFound
		}
		return dto.BookResponse{}, err
	}
	return mapper.ToBookResponse(book), nil
}

func (s *BookService) List(ctx context.Context, page params.QueryParams, category string, search string) (dto.PaginatedBookResponse, error) {
	books, err := s.BookRepo.List(ctx, page, category, search)
	if err != nil {
		return dto.PaginatedBookResponse{}, err
	}
	return mapper.ToPaginatedBookResponse(books), nil
}

func (s *BookService) Update(ctx context.Context, id string, req dto.UpdateBookRequest) (dto.BookResponse, error) {
	book, err := s.BookRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BookResponse{}, constant.ErrBookNotFound
		}
		return dto.BookResponse{}, err
	}

	// Lưu total copies cũ để tính available copies
	oldTotalCopies := book.TotalCopies

	// Cập nhật các trường
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.ISBN != nil {
		// Kiểm tra ISBN mới có trùng không (nếu thay đổi)
		if *req.ISBN != book.ISBN {
			_, err := s.BookRepo.GetByISBN(ctx, *req.ISBN)
			if err == nil {
				return dto.BookResponse{}, constant.ErrISBNAlreadyExists
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.BookResponse{}, err
			}
		}
		book.ISBN = *req.ISBN
	}
	if req.Category != nil {
		book.Category = *req.Category
	}
	if req.Description != nil {
		book.Description = req.Description
	}
	if req.TotalCopies != nil {
		newTotalCopies := *req.TotalCopies
		// Tính available copies mới
		book.AvailableCopies = book.AvailableCopies + (newTotalCopies - oldTotalCopies)
		if book.AvailableCopies < 0 {
			book.AvailableCopies = 0
		}
		book.TotalCopies = newTotalCopies
	}
	if req.PublishedYear != nil {
		book.PublishedYear = req.PublishedYear
	}
	if req.Publisher != nil {
		book.Publisher = req.Publisher
	}

	err = s.BookRepo.Update(ctx, &book)
	if err != nil {
		return dto.BookResponse{}, err
	}

	return mapper.ToBookResponse(book), nil
}

func (s *BookService) Delete(ctx context.Context, id string) error {
	// Kiểm tra book có đang được mượn không
	hasBorrows, err := s.BookRepo.HasActiveBorrows(ctx, id)
	if err != nil {
		return err
	}
	if hasBorrows {
		return constant.ErrCannotDeleteBookWithActiveBorrows
	}

	return s.BookRepo.Delete(ctx, id)
}

func (s *BookService) Search(ctx context.Context, searchReq dto.BookSearchRequest, page params.QueryParams) (dto.PaginatedBookResponse, error) {
	entitySearchReq := entity.BookSearchRequest{
		Q:             searchReq.Q,
		Category:      searchReq.Category,
		Author:        searchReq.Author,
		AvailableOnly: searchReq.AvailableOnly,
	}

	books, err := s.BookRepo.Search(ctx, entitySearchReq, page)
	if err != nil {
		return dto.PaginatedBookResponse{}, err
	}
	return mapper.ToPaginatedBookResponse(books), nil
}

