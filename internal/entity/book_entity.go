package entity

import "time"

type Book struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid"`
	Title          string    `json:"title" gorm:"not null"`
	Author         string    `json:"author" gorm:"not null"`
	ISBN           string    `json:"isbn" gorm:"not null;unique;index"`
	Category       string    `json:"category" gorm:"not null"` // Novel, Science, History, Technology
	Description    *string   `json:"description"`
	TotalCopies    int       `json:"total_copies" gorm:"default:1"`
	AvailableCopies int      `json:"available_copies" gorm:"default:1"`
	PublishedYear  *int      `json:"published_year"`
	Publisher      *string   `json:"publisher"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PaginatedBooks = Pagination[Book]

type BookSearchRequest struct {
	Q             string
	Category      string
	Author        string
	AvailableOnly bool
}

