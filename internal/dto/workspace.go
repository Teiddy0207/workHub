package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WorkspaceRequest struct {
	BaseRequest
	Name        string          `json:"name"`
	Description string          `json:"description"`
	AvatarURL   string          `json:"avatar_url"`
	IsPublic    bool            `json:"is_public"`
	Setting     json.RawMessage `json:"setting"`
}

type WorkspaceUpdateRequest struct {
	Name        *string         `json:"name" binding:"omitempty,min=1,max=255"`
	Description *string         `json:"description" binding:"omitempty,max=1000"`
	AvatarURL   *string         `json:"avatar_url"`
	IsPublic    *bool           `json:"is_public"`
	Setting     json.RawMessage `json:"setting"`
}

type WorkspaceResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	OwnerID     uuid.UUID       `json:"owner_id"`
	AvatarURL   string          `json:"avatar_url"`
	IsPublic    bool            `json:"is_public"`
	Setting     json.RawMessage `json:"setting"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type PaginatedWorkspaceResponse = Pagination[WorkspaceResponse]
