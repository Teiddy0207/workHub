package mapper

import (
	"time"
	"workHub/internal/dto"
	"workHub/internal/entity"
)

func ToBookResponse(book entity.Book) dto.BookResponse {
	return dto.BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		Author:          book.Author,
		ISBN:            book.ISBN,
		Category:        book.Category,
		Description:     book.Description,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
		PublishedYear:   book.PublishedYear,
		Publisher:       book.Publisher,
		CreatedAt:       book.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       book.UpdatedAt.Format(time.RFC3339),
	}
}

func ToPaginatedBookResponse(books entity.PaginatedBooks) dto.PaginatedBookResponse {
	var response []dto.BookResponse
	for _, book := range books.Items {
		response = append(response, ToBookResponse(book))
	}
	var totalPages int
	if books.PageSize > 0 {
		totalPages = books.TotalItems / books.PageSize
		if books.TotalItems%books.PageSize > 0 {
			totalPages++
		}
	} else {
		totalPages = 1
	}
	return dto.PaginatedBookResponse{
		Items:       response,
		TotalItems:  books.TotalItems,
		TotalPages:  totalPages,
		CurrentPage: books.PageNumber,
		PageSize:    books.PageSize,
	}
}

