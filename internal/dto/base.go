package dto

import "github.com/google/uuid"

type Pagination[T any] struct {
	Items       []T `json:"items"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

type BaseRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"-" header:"Authorization"`
}
