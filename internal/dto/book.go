package dto

type BookRequest struct {
	Title         string  `json:"title" binding:"required"`
	Author        string  `json:"author" binding:"required"`
	ISBN          string  `json:"isbn" binding:"required"`
	Category      string  `json:"category" binding:"required,oneof=Novel Science History Technology"`
	Description   *string `json:"description"`
	TotalCopies   int     `json:"total_copies" binding:"min=1"`
	PublishedYear *int    `json:"published_year"`
	Publisher     *string `json:"publisher"`
}

type BookResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Author          string  `json:"author"`
	ISBN            string  `json:"isbn"`
	Category        string  `json:"category"`
	Description     *string `json:"description"`
	TotalCopies     int     `json:"total_copies"`
	AvailableCopies int     `json:"available_copies"`
	PublishedYear   *int    `json:"published_year"`
	Publisher       *string `json:"publisher"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type UpdateBookRequest struct {
	Title         *string `json:"title"`
	Author        *string `json:"author"`
	ISBN          *string `json:"isbn"`
	Category      *string `json:"category" binding:"omitempty,oneof=Novel Science History Technology"`
	Description   *string `json:"description"`
	TotalCopies   *int    `json:"total_copies" binding:"omitempty,min=1"`
	PublishedYear *int    `json:"published_year"`
	Publisher     *string `json:"publisher"`
}

type BookSearchRequest struct {
	Q             string `form:"q"`
	Category      string `form:"category"`
	Author        string `form:"author"`
	AvailableOnly bool   `form:"available_only"`
}

type PaginatedBookResponse = Pagination[BookResponse]

