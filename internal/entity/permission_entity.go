package entity

import "time"


type Permission struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"not null;unique"`
	Action      string    `json:"action" gorm:"not null"`      
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type RolePermission struct {
	RoleID       string    `json:"role_id" gorm:"primaryKey;type:uuid;index"`
	PermissionID string    `json:"permission_id" gorm:"primaryKey;type:uuid;index"`
	CreatedAt    time.Time `json:"created_at"`
	
	Role       Role       `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserRole struct {
	UserID    string    `json:"user_id" gorm:"primaryKey;type:uuid;index"`
	RoleID    string    `json:"role_id" gorm:"primaryKey;type:uuid;index"`
	CreatedAt time.Time `json:"created_at"`
	
	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type PaginatedPermissions = Pagination[Permission]

