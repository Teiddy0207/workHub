package entity

import (
	"encoding/json"
	"time"
)

type Workspace struct {
	ID          string          `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string          `json:"name" gorm:"not null"`
	Description string          `json:"description"`
	OwnerID     string          `json:"owner_id" gorm:"not null;type:uuid;index"`
	AvatarURL   string          `json:"avatar_url"`
	IsPublic    bool            `json:"is_public" gorm:"default:false"`
	Setting     json.RawMessage `json:"setting" gorm:"type:jsonb"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type PaginatedWorkspaces = Pagination[Workspace]