package dto

import (
	// "github.com/bacabong/vnpt-be/internal/constant"
)

type Users struct {
	ID           int                 `json:"id,omitempty"`
	Name         string              `json:"name,omitempty"`
	Username     string              `json:"username,omitempty"`
	Email        string              `json:"email,omitempty"`
	Phone        string              `json:"phone,omitempty"`
	Address      string              `json:"address,omitempty"`
	Password     string              `json:"password,omitempty"`
	Role         string              `json:"role,omitempty"`
	RoleID       int                 `json:"role_id,omitempty"`
	Department   string              `json:"department,omitempty"`
	DepartmentID int                 `json:"department_id,omitempty"`
	StoreName    string              `json:"store_name,omitempty"`
	LocationID   uint                `json:"location_id,omitempty"`
	CreatedBy    int                 `json:"created_by,omitempty"`
	UpdatedBy    int                 `json:"updated_by,omitempty"`
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
	QueueID      int                 `json:"queue_id,omitempty"`
}