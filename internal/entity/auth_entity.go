package entity

import "time"

type User struct {
    ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
    Email     string    `json:"email" gorm:"not null;unique"`
    Username  string    `json:"username" gorm:"not null;unique"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedUsers = Pagination[User]