package service

import (
	"context"
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/constant"
	"workHub/logger"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

type RoleServiceInterface interface {
	CreateRole(ctx context.Context, req dto.RoleRequest) (dto.RoleResponse, error)
	GetRoleByID(ctx context.Context, id string) (dto.RoleResponse, error)
	UpdateRole(ctx context.Context, id string, req dto.RoleUpdateRequest) (dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context, params params.QueryParams) (dto.PaginatedRoleResponse, error)
}

func NewRoleService(roleRepo repository.RoleRepository) RoleServiceInterface {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, req dto.RoleRequest) (dto.RoleResponse, error) {
	logger.Info("service", "CreateRole", fmt.Sprintf("Creating new role: %s (%s)", req.Name, req.Code))

	_, err := s.roleRepo.GetRoleByCode(ctx, req.Code)
	if err == nil {
		logger.Warn("service", "CreateRole", fmt.Sprintf("Role code already exists: %s", req.Code))
		return dto.RoleResponse{}, constant.ErrTakenCredential
	}

	role := mapper.ToRoleEntity(req)

	// Tạo role trong database
	err = s.roleRepo.CreateRole(ctx, &role)
	if err != nil {
		logger.Error("service", "CreateRole", fmt.Sprintf("Failed to create role: %v", err))
		return dto.RoleResponse{}, err
	}

	createdRole, err := s.roleRepo.GetRoleByID(ctx, role.ID)
	if err != nil {
		logger.Error("service", "CreateRole", fmt.Sprintf("Failed to get created role: %v", err))
		return dto.RoleResponse{}, err
	}

	logger.Info("service", "CreateRole", fmt.Sprintf("Role created successfully: %s", createdRole.Name))
	return mapper.ToRoleResponse(createdRole), nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, id string) (dto.RoleResponse, error) {
	logger.Info("service", "GetRoleByID", fmt.Sprintf("Getting role by ID: %s", id))

	role, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		logger.Error("service", "GetRoleByID", fmt.Sprintf("Role not found: %s, error: %v", id, err))
		return dto.RoleResponse{}, constant.ErrNotFound
	}

	logger.Info("service", "GetRoleByID", fmt.Sprintf("Role found: %s", role.Name))
	return mapper.ToRoleResponse(role), nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id string, req dto.RoleUpdateRequest) (dto.RoleResponse, error) {
	logger.Info("service", "UpdateRole", fmt.Sprintf("Updating role: %s", id))

	existingRole, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		logger.Error("service", "UpdateRole", fmt.Sprintf("Role not found: %s, error: %v", id, err))
		return dto.RoleResponse{}, constant.ErrNotFound
	}

	if req.Code != nil && *req.Code != existingRole.Code {
		if _, err := s.roleRepo.GetRoleByCode(ctx, *req.Code); err == nil {
			logger.Warn("service", "UpdateRole", fmt.Sprintf("Role code already exists: %s", *req.Code))
			return dto.RoleResponse{}, constant.ErrTakenCredential
		}
	}

	updateData := mapper.ToRoleUpdateEntity(req)
	if err := s.roleRepo.UpdateRole(ctx, id, updateData); err != nil {
		logger.Error("service", "UpdateRole", fmt.Sprintf("Failed to update role: %s, error: %v", id, err))
		return dto.RoleResponse{}, err
	}

	updatedRole, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		logger.Error("service", "UpdateRole", fmt.Sprintf("Failed to get updated role: %s, error: %v", id, err))
		return dto.RoleResponse{}, err
	}

	logger.Info("service", "UpdateRole", fmt.Sprintf("Role updated successfully: %s", updatedRole.Name))
	return mapper.ToRoleResponse(updatedRole), nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	logger.Info("service", "DeleteRole", fmt.Sprintf("Deleting role: %s", id))

	_, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		logger.Error("service", "DeleteRole", fmt.Sprintf("Role not found: %s, error: %v", id, err))
		return constant.ErrNotFound
	}

	err = s.roleRepo.DeleteRole(ctx, id)
	if err != nil {
		logger.Error("service", "DeleteRole", fmt.Sprintf("Failed to delete role: %s, error: %v", id, err))
		return err
	}

	logger.Info("service", "DeleteRole", fmt.Sprintf("Role deleted successfully: %s", id))
	return nil
}

func (s *RoleService) ListRoles(ctx context.Context, params params.QueryParams) (dto.PaginatedRoleResponse, error) {
	logger.Info("service", "ListRoles", fmt.Sprintf("Listing roles with params: page=%d, size=%d, search=%s", 
		params.PageNumber, params.PageSize, params.Search))

	roles, err := s.roleRepo.ListRoles(ctx, params)
	if err != nil {
		logger.Error("service", "ListRoles", fmt.Sprintf("Failed to list roles: %v", err))
		return dto.PaginatedRoleResponse{}, err
	}

	logger.Info("service", "ListRoles", fmt.Sprintf("Listed %d roles successfully", len(roles.Items)))
	return mapper.ToPaginatedRoleResponse(roles), nil
}