package dto

type BorrowRequest struct {
	BookID        string `json:"book_id" binding:"required"`
	DaysToBorrow  int    `json:"days_to_borrow" binding:"required,min=1,max=30"`
}

type BorrowResponse struct {
	ID         string       `json:"id"`
	UserID     string       `json:"user_id"`
	BookID     string       `json:"book_id"`
	BorrowDate string       `json:"borrow_date"`
	DueDate    string       `json:"due_date"`
	ReturnDate *string      `json:"return_date"`
	Status     string       `json:"status"`
	Notes      *string      `json:"notes"`
	CreatedAt  string       `json:"created_at"`
	UpdatedAt  string       `json:"updated_at"`
	User       *UserInfo    `json:"user,omitempty"`
	Book       *BookResponse `json:"book,omitempty"`
}

type PaginatedBorrowResponse = Pagination[BorrowResponse]

type ReturnBorrowRequest struct {
	Notes *string `json:"notes"`
}

