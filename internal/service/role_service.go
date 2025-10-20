package service

import (
	"context"
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/entity"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/constant"
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
	fmt.Printf("🎯 Creating new role: %s (%s)\n", req.Name, req.Code)

	// Kiểm tra xem role code đã tồn tại chưa
	_, err := s.roleRepo.GetRoleByCode(ctx, req.Code)
	if err == nil {
		fmt.Printf("❌ Role code already exists: %s\n", req.Code)
		return dto.RoleResponse{}, constant.ErrTakenCredential
	}

	// Chuyển đổi DTO thành entity
	role := mapper.ToRoleEntity(req)

	// Tạo role trong database
	err = s.roleRepo.CreateRole(ctx, &role)
	if err != nil {
		fmt.Printf("❌ Failed to create role: %v\n", err)
		return dto.RoleResponse{}, err
	}

	// Lấy role vừa tạo để trả về response
	createdRole, err := s.roleRepo.GetRoleByCode(ctx, req.Code)
	if err != nil {
		fmt.Printf("❌ Failed to get created role: %v\n", err)
		return dto.RoleResponse{}, err
	}

	fmt.Printf("✅ Role created successfully: %s\n", createdRole.Name)
	return mapper.ToRoleResponse(createdRole), nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, id string) (dto.RoleResponse, error) {
	fmt.Printf("🎯 Getting role by ID: %s\n", id)

	role, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		fmt.Printf("❌ Role not found: %v\n", err)
		return dto.RoleResponse{}, constant.ErrNotFound
	}

	fmt.Printf("✅ Role found: %s\n", role.Name)
	return mapper.ToRoleResponse(role), nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id string, req dto.RoleUpdateRequest) (dto.RoleResponse, error) {
	fmt.Printf("🎯 Updating role: %s\n", id)

	// Kiểm tra role có tồn tại không
	existingRole, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		fmt.Printf("❌ Role not found: %v\n", err)
		return dto.RoleResponse{}, constant.ErrNotFound
	}

	// Cập nhật các field được cung cấp
	updateData := &entity.Role{}
	
	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Code != nil {
		// Kiểm tra xem code mới có trùng với role khác không
		if *req.Code != existingRole.Code {
			_, err := s.roleRepo.GetRoleByCode(ctx, *req.Code)
			if err == nil {
				fmt.Printf("❌ Role code already exists: %s\n", *req.Code)
				return dto.RoleResponse{}, constant.ErrTakenCredential
			}
		}
		updateData.Code = *req.Code
	}
	if req.Description != nil {
		updateData.Description = *req.Description
	}
	if req.IsActive != nil {
		updateData.IsActive = *req.IsActive
	}

	// Cập nhật role
	err = s.roleRepo.UpdateRole(ctx, id, updateData)
	if err != nil {
		fmt.Printf("❌ Failed to update role: %v\n", err)
		return dto.RoleResponse{}, err
	}

	// Lấy role đã cập nhật
	updatedRole, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		fmt.Printf("❌ Failed to get updated role: %v\n", err)
		return dto.RoleResponse{}, err
	}

	fmt.Printf("✅ Role updated successfully: %s\n", updatedRole.Name)
	return mapper.ToRoleResponse(updatedRole), nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	fmt.Printf("🎯 Deleting role: %s\n", id)

	// Kiểm tra role có tồn tại không
	_, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		fmt.Printf("❌ Role not found: %v\n", err)
		return constant.ErrNotFound
	}

	// Xóa role (soft delete)
	err = s.roleRepo.DeleteRole(ctx, id)
	if err != nil {
		fmt.Printf("❌ Failed to delete role: %v\n", err)
		return err
	}

	fmt.Printf("✅ Role deleted successfully: %s\n", id)
	return nil
}

func (s *RoleService) ListRoles(ctx context.Context, params params.QueryParams) (dto.PaginatedRoleResponse, error) {
	fmt.Printf("🎯 Listing roles with params: page=%d, size=%d, search=%s\n", 
		params.PageNumber, params.PageSize, params.Search)

	roles, err := s.roleRepo.ListRoles(ctx, params)
	if err != nil {
		fmt.Printf("❌ Failed to list roles: %v\n", err)
		return dto.PaginatedRoleResponse{}, err
	}

	fmt.Printf("✅ Listed %d roles successfully\n", len(roles.Items))
	return mapper.ToPaginatedRoleResponse(roles), nil
}