package dto

import (
	// "github.com/bacabong/vnpt-be/internal/constant"
)

type Users struct {
	ID           string              `json:"id,omitempty"` // UUID string
	Name         string              `json:"name,omitempty"`
	Username     string              `json:"username,omitempty"`
	Email        string              `json:"email,omitempty"`
	Phone        string              `json:"phone,omitempty"`
	Address      string              `json:"address,omitempty"`
	Password     string              `json:"password,omitempty"`
	Role         string              `json:"role,omitempty"`
	RoleID       string              `json:"role_id,omitempty"` // UUID string
	Department   string              `json:"department,omitempty"`
	DepartmentID string              `json:"department_id,omitempty"` // UUID string
	StoreName    string              `json:"store_name,omitempty"`
	LocationID   string              `json:"location_id,omitempty"` // UUID string
	CreatedBy    string              `json:"created_by,omitempty"` // UUID string
	UpdatedBy    string              `json:"updated_by,omitempty"` // UUID string
	AvatarURL    string              `json:"avatar_url,omitempty"`
	Avatar       string              `json:"avatar,omitempty"`
	// Gender       constant.Gender     `json:"gender,omitempty"`
	// DOB          time.Time           `json:"dob,omitempty"`
	Bio          string              `json:"bio,omitempty"`
	AuthTwoFace  string              `json:"auth_two_face,omitempty"`
	IsVIP        bool                `json:"is_vip,omitempty"`
	// CreatedAt    time.Time           `json:"created_at,omitempty"`
	// UpdatedAt    time.Time           `json:"updated_at,omitempty"`
	// DeletedAt    *time.Time          `json:"deleted_at,omitempty"`
	QueueID      string              `json:"queue_id,omitempty"` // UUID string
}