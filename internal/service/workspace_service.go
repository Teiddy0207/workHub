package service

import (
	"context"
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/logger"
)

type workService struct {
	workspaceRepo repository.WorkspaceRepository
}

type WorkspaceServiceInterface interface {
	CreateWorkspace(ctx context.Context, req dto.WorkspaceRequest) (dto.WorkspaceResponse, error)
}

func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository) WorkspaceServiceInterface {
	return &workService{
		workspaceRepo: workspaceRepo,
	}
}

func (w *workService) CreateWorkspace(ctx context.Context, req dto.WorkspaceRequest) (dto.WorkspaceResponse, error) {
	logger.Info("service", "CreateWorkspace", fmt.Sprintf("Creating new workspace: %s", req.Name))
	workspace := mapper.ToWorkspaceEntity(req)
	err := w.workspaceRepo.CreateWorkspace(ctx, &workspace)
	if err != nil {
		logger.Error("service", "CreateWorkspace", fmt.Sprintf("Failed to create workspace: %v", err))
		return dto.WorkspaceResponse{}, err
	}
	createdWorkspace, err := w.workspaceRepo.GetWorkspaceByID(ctx, workspace.ID)
	if err != nil {
		logger.Error("service", "CreateWorkspace", fmt.Sprintf("Failed to get created workspace: %v", err))
		return dto.WorkspaceResponse{}, err
	}
	logger.Info("service", "CreateWorkspace", fmt.Sprintf("Workspace created successfully: %s", createdWorkspace.Name))
	return mapper.ToWorkspaceResponse(createdWorkspace), nil
}
