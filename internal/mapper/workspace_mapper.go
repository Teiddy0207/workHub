package mapper

import (
	"time"
	"workHub/internal/dto"
	"workHub/internal/entity"

	"github.com/google/uuid"
)

// ToWorkspaceEntity chuyển đổi dto.WorkspaceRequest thành entity.Workspace
func ToWorkspaceEntity(req dto.WorkspaceRequest) entity.Workspace {
	now := time.Now()

	// Set default values nếu cần
	setting := req.Setting
	if len(setting) == 0 {
		setting = []byte("{}")
	}

	return entity.Workspace{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID.String(),
		AvatarURL:   req.AvatarURL,
		IsPublic:    req.IsPublic,
		Setting:     setting,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ToWorkspaceResponse chuyển đổi entity.Workspace thành dto.WorkspaceResponse
func ToWorkspaceResponse(workspace entity.Workspace) dto.WorkspaceResponse {
	ownerID, _ := uuid.Parse(workspace.OwnerID)
	workspaceID, _ := uuid.Parse(workspace.ID)

	return dto.WorkspaceResponse{
		ID:          workspaceID,
		Name:        workspace.Name,
		Description: workspace.Description,
		OwnerID:     ownerID,
		AvatarURL:   workspace.AvatarURL,
		IsPublic:    workspace.IsPublic,
		Setting:     workspace.Setting,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
}
