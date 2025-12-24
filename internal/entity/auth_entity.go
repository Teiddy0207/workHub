package entity

import "time"

type User struct {
    ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
    Email       string    `json:"email" gorm:"not null;unique;index"`
    Username    string    `json:"username" gorm:"not null;unique"`
    Password    string    `json:"-" gorm:"not null"`
    FullName    string    `json:"full_name" gorm:"not null"`
    PhoneNumber *string   `json:"phone_number"`
    Address     *string   `json:"address"`
    Role        string    `json:"role" gorm:"not null;default:'student'"` // student, teacher, admin
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}



type PaginatedUsers = Pagination[User]