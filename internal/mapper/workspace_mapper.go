package mapper

import (
	"time"
	"workHub/internal/dto"
	"workHub/internal/entity"

	"github.com/google/uuid"
)

func ToWorkspaceEntity(req dto.WorkspaceRequest) entity.Workspace {
	now := time.Now()

	setting := req.Setting
	if len(setting) == 0 {
		setting = []byte("{}")
	}

	return entity.Workspace{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.UserID.String(),
		AvatarURL:   req.AvatarURL,
		IsPublic:    req.IsPublic,
		Setting:     setting,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

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

func ToWorkspaceResponseList(workspaces []entity.Workspace) []dto.WorkspaceResponse {
	responses := make([]dto.WorkspaceResponse, len(workspaces))
	for i, workspace := range workspaces {
		responses[i] = ToWorkspaceResponse(workspace)
	}
	return responses
}

func ToPaginatedWorkspaceResponse(paginatedWorkspaces entity.PaginatedWorkspaces) dto.PaginatedWorkspaceResponse {
	return dto.PaginatedWorkspaceResponse{
		Items:       ToWorkspaceResponseList(paginatedWorkspaces.Items),
		TotalItems:  paginatedWorkspaces.TotalItems,
		TotalPages:  (paginatedWorkspaces.TotalItems + paginatedWorkspaces.PageSize - 1) / paginatedWorkspaces.PageSize,
		CurrentPage: paginatedWorkspaces.PageNumber,
		PageSize:    paginatedWorkspaces.PageSize,
	}
}

func ToWorkspaceUpdateEntity(req dto.WorkspaceUpdateRequest) *entity.Workspace {
	updateData := &entity.Workspace{}

	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Description != nil {
		updateData.Description = *req.Description
	}
	if req.AvatarURL != nil {
		updateData.AvatarURL = *req.AvatarURL
	}
	if req.IsPublic != nil {
		updateData.IsPublic = *req.IsPublic
	}
	if len(req.Setting) > 0 {
		updateData.Setting = req.Setting
	}

	return updateData
}
