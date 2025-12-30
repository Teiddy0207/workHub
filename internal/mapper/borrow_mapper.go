package mapper

import (
	"time"
	"workHub/internal/dto"
	"workHub/internal/entity"
)

func ToBorrowResponse(borrow entity.Borrow, includeUser bool, includeBook bool) dto.BorrowResponse {
	resp := dto.BorrowResponse{
		ID:         borrow.ID,
		UserID:     borrow.UserID,
		BookID:     borrow.BookID,
		BorrowDate: borrow.BorrowDate.Format("2006-01-02"),
		DueDate:    borrow.DueDate.Format("2006-01-02"),
		Status:     borrow.Status,
		Notes:      borrow.Notes,
		CreatedAt:  borrow.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  borrow.UpdatedAt.Format(time.RFC3339),
	}

	if borrow.ReturnDate != nil {
		returnDate := borrow.ReturnDate.Format("2006-01-02")
		resp.ReturnDate = &returnDate
	}

	if includeUser {
		userInfo := dto.UserInfo{
			ID:          borrow.User.ID,
			Email:       borrow.User.Email,
			FullName:    borrow.User.FullName,
			PhoneNumber: borrow.User.PhoneNumber,
			Address:     borrow.User.Address,
			Role:        borrow.User.Role,
		}
		resp.User = &userInfo
	}

	if includeBook {
		bookResp := ToBookResponse(borrow.Book)
		resp.Book = &bookResp
	}

	return resp
}

func ToPaginatedBorrowResponse(borrows entity.PaginatedBorrows, includeUser bool, includeBook bool) dto.PaginatedBorrowResponse {
	var response []dto.BorrowResponse
	for _, borrow := range borrows.Items {
		response = append(response, ToBorrowResponse(borrow, includeUser, includeBook))
	}
	var totalPages int
	if borrows.PageSize > 0 {
		totalPages = borrows.TotalItems / borrows.PageSize
		if borrows.TotalItems%borrows.PageSize > 0 {
			totalPages++
		}
	} else {
		totalPages = 1
	}
	return dto.PaginatedBorrowResponse{
		Items:       response,
		TotalItems:  borrows.TotalItems,
		TotalPages:  totalPages,
		CurrentPage: borrows.PageNumber,
		PageSize:    borrows.PageSize,
	}
}

