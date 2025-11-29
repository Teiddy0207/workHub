package service

import (
	"context"
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/constant"
	"workHub/logger"

	"github.com/google/uuid"
)

type workService struct {
	workspaceRepo repository.WorkspaceRepository
	authRepo      repository.AuthRepository
}

type WorkspaceServiceInterface interface {
	CreateWorkspace(ctx context.Context, req dto.WorkspaceRequest, userEmail string) (dto.WorkspaceResponse, error)
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
