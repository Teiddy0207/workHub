package repository

import (
	"context"
	"fmt"
	"workHub/internal/entity"
	"workHub/logger"

	"gorm.io/gorm"
)

type workspaceRepository struct {
	db *gorm.DB
}

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, workspace *entity.Workspace) error
	GetWorkspaceByID(ctx context.Context, id string) (entity.Workspace, error)
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) CreateWorkspace(ctx context.Context, workspace *entity.Workspace) error {
	logger.Info("repository", "CreateWorkspace", fmt.Sprintf("Creating workspace: %s", workspace.Name))

	err := r.db.WithContext(ctx).Create(workspace).Error
	if err != nil {
		logger.Error("repository", "CreateWorkspace", fmt.Sprintf("Failed to create workspace: %v", err))
		return err
	}

	logger.Info("repository", "CreateWorkspace", fmt.Sprintf("Workspace created successfully: %s", workspace.ID))
	return nil
}

func (r *workspaceRepository) GetWorkspaceByID(ctx context.Context, id string) (entity.Workspace, error) {
	logger.Info("repository", "GetWorkspaceByID", fmt.Sprintf("Getting workspace by ID: %s", id))

	var workspace entity.Workspace
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&workspace).Error

	if err != nil {
		logger.Error("repository", "GetWorkspaceByID", fmt.Sprintf("Workspace not found: %s, error: %v", id, err))
		return entity.Workspace{}, err
	}

	logger.Info("repository", "GetWorkspaceByID", fmt.Sprintf("Workspace found: %s", workspace.Name))
	return workspace, nil
}

