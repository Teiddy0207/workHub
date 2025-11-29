package service

import (
	"context"
	"fmt"
	"workHub/constant"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/logger"
	"workHub/pkg/params"

	"github.com/google/uuid"
)

type workService struct {
	workspaceRepo repository.WorkspaceRepository
	authRepo      repository.AuthRepository
}

type WorkspaceServiceInterface interface {
	CreateWorkspace(ctx context.Context, req dto.WorkspaceRequest, userEmail string) (dto.WorkspaceResponse, error)
	GetWorkspaceByID(ctx context.Context, id string) (dto.WorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, id string, req dto.WorkspaceUpdateRequest) (dto.WorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, id string) error
	ListWorkspaces(ctx context.Context, params params.QueryParams) (dto.PaginatedWorkspaceResponse, error)
}

func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository, authRepo repository.AuthRepository) WorkspaceServiceInterface {
	return &workService{
		workspaceRepo: workspaceRepo,
		authRepo:      authRepo,
	}
}

func (w *workService) CreateWorkspace(ctx context.Context, req dto.WorkspaceRequest, userEmail string) (dto.WorkspaceResponse, error) {
	logger.Info("service", "CreateWorkspace", fmt.Sprintf("Creating new workspace: %s for user: %s", req.Name, userEmail))

	// Query database để lấy user UUID từ email
	user, err := w.authRepo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		logger.Error("service", "CreateWorkspace", fmt.Sprintf("Failed to get user by email: %v", err))
		return dto.WorkspaceResponse{}, constant.ErrNotFound
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		logger.Error("service", "CreateWorkspace", fmt.Sprintf("Invalid user ID format: %v", err))
		return dto.WorkspaceResponse{}, constant.ErrInvalidInput
	}

	req.UserID = userID

	workspace := mapper.ToWorkspaceEntity(req)
	err = w.workspaceRepo.CreateWorkspace(ctx, &workspace)
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

func (w *workService) GetWorkspaceByID(ctx context.Context, id string) (dto.WorkspaceResponse, error) {
	logger.Info("service", "GetWorkspaceByID", fmt.Sprintf("Getting workspace by ID: %s", id))

	workspace, err := w.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		logger.Error("service", "GetWorkspaceByID", fmt.Sprintf("Workspace not found: %s, error: %v", id, err))
		return dto.WorkspaceResponse{}, constant.ErrNotFound
	}

	logger.Info("service", "GetWorkspaceByID", fmt.Sprintf("Workspace found: %s", workspace.Name))
	return mapper.ToWorkspaceResponse(workspace), nil
}

func (w *workService) UpdateWorkspace(ctx context.Context, id string, req dto.WorkspaceUpdateRequest) (dto.WorkspaceResponse, error) {
	logger.Info("service", "UpdateWorkspace", fmt.Sprintf("Updating workspace: %s", id))

	_, err := w.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		logger.Error("service", "UpdateWorkspace", fmt.Sprintf("Workspace not found: %s, error: %v", id, err))
		return dto.WorkspaceResponse{}, constant.ErrNotFound
	}

	updateData := mapper.ToWorkspaceUpdateEntity(req)
	if err := w.workspaceRepo.UpdateWorkspace(ctx, id, updateData); err != nil {
		logger.Error("service", "UpdateWorkspace", fmt.Sprintf("Failed to update workspace: %s, error: %v", id, err))
		return dto.WorkspaceResponse{}, err
	}

	updatedWorkspace, err := w.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		logger.Error("service", "UpdateWorkspace", fmt.Sprintf("Failed to get updated workspace: %s, error: %v", id, err))
		return dto.WorkspaceResponse{}, err
	}

	logger.Info("service", "UpdateWorkspace", fmt.Sprintf("Workspace updated successfully: %s", updatedWorkspace.Name))
	return mapper.ToWorkspaceResponse(updatedWorkspace), nil
}

func (w *workService) DeleteWorkspace(ctx context.Context, id string) error {
	logger.Info("service", "DeleteWorkspace", fmt.Sprintf("Deleting workspace: %s", id))

	_, err := w.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		logger.Error("service", "DeleteWorkspace", fmt.Sprintf("Workspace not found: %s, error: %v", id, err))
		return constant.ErrNotFound
	}

	err = w.workspaceRepo.DeleteWorkspace(ctx, id)
	if err != nil {
		logger.Error("service", "DeleteWorkspace", fmt.Sprintf("Failed to delete workspace: %s, error: %v", id, err))
		return err
	}

	logger.Info("service", "DeleteWorkspace", fmt.Sprintf("Workspace deleted successfully: %s", id))
	return nil
}

func (w *workService) ListWorkspaces(ctx context.Context, params params.QueryParams) (dto.PaginatedWorkspaceResponse, error) {
	logger.Info("service", "ListWorkspaces", fmt.Sprintf("Listing workspaces with params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	workspaces, err := w.workspaceRepo.ListWorkspaces(ctx, params)
	if err != nil {
		logger.Error("service", "ListWorkspaces", fmt.Sprintf("Failed to list workspaces: %v", err))
		return dto.PaginatedWorkspaceResponse{}, err
	}

	logger.Info("service", "ListWorkspaces", fmt.Sprintf("Listed %d workspaces successfully", len(workspaces.Items)))
	return mapper.ToPaginatedWorkspaceResponse(workspaces), nil
}
