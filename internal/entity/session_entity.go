package entity

import "time"

// Session - Entity cho bảng sessions để quản lý phiên người dùng
type Session struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid"`
	UserID       string    `json:"user_id" gorm:"not null;type:uuid;index"`
	AccessToken  string    `json:"access_token" gorm:"type:text;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text;not null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null;index"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

