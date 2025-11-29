package repository

import (
	"context"
	"fmt"
	"workHub/internal/entity"
	"workHub/logger"
	"workHub/pkg/params"

	"gorm.io/gorm"
)

type workspaceRepository struct {
	db *gorm.DB
}

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, workspace *entity.Workspace) error
	GetWorkspaceByID(ctx context.Context, id string) (entity.Workspace, error)
	UpdateWorkspace(ctx context.Context, id string, workspace *entity.Workspace) error
	DeleteWorkspace(ctx context.Context, id string) error
	ListWorkspaces(ctx context.Context, params params.QueryParams) (entity.PaginatedWorkspaces, error)
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

func (r *workspaceRepository) UpdateWorkspace(ctx context.Context, id string, workspace *entity.Workspace) error {
	logger.Info("repository", "UpdateWorkspace", fmt.Sprintf("Updating workspace: %s", id))

	result := r.db.WithContext(ctx).
		Model(&entity.Workspace{}).
		Where("id = ?", id).
		Updates(workspace)

	if result.Error != nil {
		logger.Error("repository", "UpdateWorkspace", fmt.Sprintf("Failed to update workspace: %s, error: %v", id, result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Error("repository", "UpdateWorkspace", fmt.Sprintf("Workspace not found for update: %s", id))
		return gorm.ErrRecordNotFound
	}

	logger.Info("repository", "UpdateWorkspace", fmt.Sprintf("Workspace updated successfully: %s", id))
	return nil
}

func (r *workspaceRepository) DeleteWorkspace(ctx context.Context, id string) error {
	logger.Info("repository", "DeleteWorkspace", fmt.Sprintf("Hard deleting workspace: %s", id))

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.Workspace{})

	if result.Error != nil {
		logger.Error("repository", "DeleteWorkspace", fmt.Sprintf("Failed to delete workspace: %s, error: %v", id, result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Error("repository", "DeleteWorkspace", fmt.Sprintf("Workspace not found for deletion: %s", id))
		return gorm.ErrRecordNotFound
	}

	logger.Info("repository", "DeleteWorkspace", fmt.Sprintf("Workspace deleted successfully: %s", id))
	return nil
}

func (r *workspaceRepository) ListWorkspaces(ctx context.Context, params params.QueryParams) (entity.PaginatedWorkspaces, error) {
	logger.Info("repository", "ListWorkspaces", fmt.Sprintf("Listing workspaces with params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	var workspaces []entity.Workspace
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Workspace{})

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			searchTerm, searchTerm)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		logger.Error("repository", "ListWorkspaces", fmt.Sprintf("Failed to count workspaces: %v", err))
		return entity.PaginatedWorkspaces{}, err
	}

	offset := (params.PageNumber - 1) * params.PageSize

	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&workspaces).Error; err != nil {
		logger.Error("repository", "ListWorkspaces", fmt.Sprintf("Failed to fetch workspaces: %v", err))
		return entity.PaginatedWorkspaces{}, err
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	logger.Info("repository", "ListWorkspaces", fmt.Sprintf("Found %d workspaces (total: %d)", len(workspaces), totalItems))

	return entity.PaginatedWorkspaces{
		Items:      workspaces,
		TotalItems: int(totalItems),
		PageNumber: params.PageNumber,
		PageSize:   pageSize,
	}, nil
}

