package entity

import (
	"encoding/json"
	"time"
)

type WorkspaceMember struct {
	ID                   string          `json:"id" gorm:"primaryKey;type:uuid"`
	WorkspaceID          string          `json:"workspace_id" gorm:"not null;type:uuid;index"`
	UserID               string          `json:"user_id" gorm:"not null;type:uuid;index"`
	RoleID               string          `json:"role_id" gorm:"not null;type:uuid;index"`
	JoinedAt             time.Time       `json:"joined_at" gorm:"not null;default:now()"`
	InvitedBy            *string         `json:"invited_by" gorm:"type:uuid;index"`
	IsActive             bool            `json:"is_active" gorm:"default:true"`
	LastSeenAt           *time.Time      `json:"last_seen_at"`
	NotificationSettings json.RawMessage `json:"notification_settings" gorm:"type:jsonb;default:'{}'"`
	CustomStatus         *string         `json:"custom_status" gorm:"type:varchar(100)"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`

	Workspace Workspace `json:"workspace" gorm:"foreignKey:WorkspaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Inviter   *User     `json:"inviter" gorm:"foreignKey:InvitedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type PaginatedWorkspaceMembers = Pagination[WorkspaceMember]

