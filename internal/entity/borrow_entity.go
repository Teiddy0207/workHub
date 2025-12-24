package entity

import "time"

type Borrow struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid"`
	UserID     string    `json:"user_id" gorm:"not null;type:uuid;index"`
	BookID     string    `json:"book_id" gorm:"not null;type:uuid;index"`
	BorrowDate time.Time `json:"borrow_date" gorm:"not null;type:date"`
	DueDate    time.Time `json:"due_date" gorm:"not null;type:date"`
	ReturnDate *time.Time `json:"return_date" gorm:"type:date"`
	Status     string    `json:"status" gorm:"not null;default:'pending'"` // pending, borrowed, returned, overdue
	Notes      *string   `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Book Book `json:"book" gorm:"foreignKey:BookID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type PaginatedBorrows = Pagination[Borrow]

