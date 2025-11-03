package controller

import (
	"errors"
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/pkg/handler"
	"workHub/pkg/params"
	"workHub/constant"
	"workHub/logger"
	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	handler.BaseHandler
	service service.PermissionServiceInterface
}

func NewPermissionController(service service.PermissionServiceInterface) *PermissionController {
	return &PermissionController{
		BaseHandler: handler.NewBaseHandler(),
		service:     service,
	}
}

// CreatePermission - Tạo permission mới
func (p *PermissionController) CreatePermission(c *gin.Context) {
	logger.Info("controller", "CreatePermission", "CreatePermission controller called")
	ctx := c.Request.Context()

	var req dto.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "CreatePermission", fmt.Sprintf("Bind JSON error: %v", err))
		p.BaseHandler.BadRequest(c, constant.PERMISSION_CREATED_FAIL)
		return
	}

	logger.Info("controller", "CreatePermission", fmt.Sprintf("Request received: name=%s, code=%s", req.Name, req.Code))

	response, err := p.service.CreatePermission(ctx, req)
	if err != nil {
		logger.Error("controller", "CreatePermission", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrTakenCredential) || errors.Is(err, constant.ErrNotFound) {
			p.BaseHandler.BadRequest(c, err.Error())
		} else {
			p.BaseHandler.BadRequest(c, constant.PERMISSION_CREATED_FAIL)
		}
		return
	}

	logger.Info("controller", "CreatePermission", "Permission created successfully")
	p.BaseHandler.SuccessResponse(c, response, constant.PERMISSION_CREATED_SUCCESSFULLY)
}

// GetPermissionByID - Lấy thông tin permission theo ID
func (p *PermissionController) GetPermissionByID(c *gin.Context) {
	logger.Info("controller", "GetPermissionByID", "GetPermissionByID controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "GetPermissionByID", "Missing permission ID")
		p.BaseHandler.BadRequest(c, "permission ID is required")
		return
	}

	logger.Info("controller", "GetPermissionByID", fmt.Sprintf("Getting permission with ID: %s", id))

	response, err := p.service.GetPermissionByID(ctx, id)
	if err != nil {
		logger.Error("controller", "GetPermissionByID", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, err.Error())
		return
	}

	logger.Info("controller", "GetPermissionByID", "Permission retrieved successfully")
	p.BaseHandler.SuccessResponse(c, response, constant.PERMISSION_GET_SUCCESSFULLY)
}

// UpdatePermission - Cập nhật permission
func (p *PermissionController) UpdatePermission(c *gin.Context) {
	logger.Info("controller", "UpdatePermission", "UpdatePermission controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "UpdatePermission", "Missing permission ID")
		p.BaseHandler.BadRequest(c, "permission ID is required")
		return
	}

	var req dto.PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "UpdatePermission", fmt.Sprintf("Bind JSON error: %v", err))
		p.BaseHandler.BadRequest(c, constant.PERMISSION_UPDATED_FAIL)
		return
	}

	logger.Info("controller", "UpdatePermission", fmt.Sprintf("Updating permission with ID: %s", id))

	response, err := p.service.UpdatePermission(ctx, id, req)
	if err != nil {
		logger.Error("controller", "UpdatePermission", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) || errors.Is(err, constant.ErrTakenCredential) {
			p.BaseHandler.BadRequest(c, err.Error())
		} else {
			p.BaseHandler.BadRequest(c, constant.PERMISSION_UPDATED_FAIL)
		}
		return
	}

	logger.Info("controller", "UpdatePermission", "Permission updated successfully")
	p.BaseHandler.SuccessResponse(c, response, constant.PERMISSION_UPDATED_SUCCESSFULLY)
}

// DeletePermission - Xóa permission
func (p *PermissionController) DeletePermission(c *gin.Context) {
	logger.Info("controller", "DeletePermission", "DeletePermission controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "DeletePermission", "Missing permission ID")
		p.BaseHandler.BadRequest(c, "permission ID is required")
		return
	}

	logger.Info("controller", "DeletePermission", fmt.Sprintf("Deleting permission with ID: %s", id))

	err := p.service.DeletePermission(ctx, id)
	if err != nil {
		logger.Error("controller", "DeletePermission", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) {
			p.BaseHandler.BadRequest(c, err.Error())
		} else {
			p.BaseHandler.BadRequest(c, constant.PERMISSION_DELETED_FAIL)
		}
		return
	}

	logger.Info("controller", "DeletePermission", "Permission deleted successfully")
	p.BaseHandler.SuccessResponse(c, nil, constant.PERMISSION_DELETED_SUCCESSFULLY)
}

// ListPermissions - Lấy danh sách permissions
func (p *PermissionController) ListPermissions(c *gin.Context) {
	logger.Info("controller", "ListPermissions", "ListPermissions controller called")
	ctx := c.Request.Context()

	params := params.NewQueryParams(c)
	logger.Info("controller", "ListPermissions", fmt.Sprintf("Query params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	response, err := p.service.ListPermissions(ctx, params)
	if err != nil {
		logger.Error("controller", "ListPermissions", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, constant.PERMISSION_GET_LIST_FAIL)
		return
	}

	logger.Info("controller", "ListPermissions", "Permissions listed successfully")
	p.BaseHandler.SuccessResponse(c, response, constant.PERMISSION_GET_LIST_SUCCESSFULLY)
}

// AssignPermissionsToRole - Gán permissions cho role
func (p *PermissionController) AssignPermissionsToRole(c *gin.Context) {
	logger.Info("controller", "AssignPermissionsToRole", "AssignPermissionsToRole controller called")
	ctx := c.Request.Context()

	roleID := c.Param("id")
	if roleID == "" {
		logger.Warn("controller", "AssignPermissionsToRole", "Missing role ID")
		p.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	var req dto.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "AssignPermissionsToRole", fmt.Sprintf("Bind JSON error: %v", err))
		p.BaseHandler.BadRequest(c, constant.ASSIGN_PERMISSION_FAIL)
		return
	}

	logger.Info("controller", "AssignPermissionsToRole", fmt.Sprintf("Assigning %d permissions to role: %s", len(req.PermissionIDs), roleID))

	err := p.service.AssignPermissionsToRole(ctx, roleID, req)
	if err != nil {
		logger.Error("controller", "AssignPermissionsToRole", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, constant.ASSIGN_PERMISSION_FAIL)
		return
	}

	logger.Info("controller", "AssignPermissionsToRole", "Permissions assigned successfully")
	p.BaseHandler.SuccessResponse(c, nil, constant.ASSIGN_PERMISSION)
}

// RemovePermissionsFromRole - Xóa permissions khỏi role
func (p *PermissionController) RemovePermissionsFromRole(c *gin.Context) {
	logger.Info("controller", "RemovePermissionsFromRole", "RemovePermissionsFromRole controller called")
	ctx := c.Request.Context()

	roleID := c.Param("id")
	if roleID == "" {
		logger.Warn("controller", "RemovePermissionsFromRole", "Missing role ID")
		p.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	var req dto.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "RemovePermissionsFromRole", fmt.Sprintf("Bind JSON error: %v", err))
		p.BaseHandler.BadRequest(c, constant.ASSIGN_PERMISSION_FAIL)
		return
	}

	logger.Info("controller", "RemovePermissionsFromRole", fmt.Sprintf("Removing %d permissions from role: %s", len(req.PermissionIDs), roleID))

	err := p.service.RemovePermissionsFromRole(ctx, roleID, req)
	if err != nil {
		logger.Error("controller", "RemovePermissionsFromRole", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, constant.ASSIGN_PERMISSION_FAIL)
		return
	}

	logger.Info("controller", "RemovePermissionsFromRole", "Permissions removed successfully")
	p.BaseHandler.SuccessResponse(c, nil, "Xóa quyền khỏi role thành công")
}

// GetRoleWithPermissions - Lấy role kèm danh sách permissions có phân trang
func (p *PermissionController) GetRoleWithPermissions(c *gin.Context) {
	logger.Info("controller", "GetRoleWithPermissions", "GetRoleWithPermissions controller called")
	ctx := c.Request.Context()

	roleID := c.Param("id")
	if roleID == "" {
		logger.Warn("controller", "GetRoleWithPermissions", "Missing role ID")
		p.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	params := params.NewQueryParams(c)
	logger.Info("controller", "GetRoleWithPermissions", fmt.Sprintf("Query params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	response, err := p.service.GetRoleWithPermissions(ctx, roleID, params)
	if err != nil {
		logger.Error("controller", "GetRoleWithPermissions", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, err.Error())
		return
	}

	logger.Info("controller", "GetRoleWithPermissions", "Role with permissions retrieved successfully")
	p.BaseHandler.SuccessResponse(c, response, "Lấy role và permissions thành công")
}

// AssignRolesToUser - Gán roles cho user
func (p *PermissionController) AssignRolesToUser(c *gin.Context) {
	logger.Info("controller", "AssignRolesToUser", "AssignRolesToUser controller called")
	ctx := c.Request.Context()

	userID := c.Param("id")
	if userID == "" {
		logger.Warn("controller", "AssignRolesToUser", "Missing user ID")
		p.BaseHandler.BadRequest(c, "user ID is required")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "AssignRolesToUser", fmt.Sprintf("Bind JSON error: %v", err))
		p.BaseHandler.BadRequest(c, "Gán role cho user thất bại")
		return
	}

	err := p.service.AssignRolesToUser(ctx, userID, req)
	if err != nil {
		logger.Error("controller", "AssignRolesToUser", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, "Gán role cho user thất bại")
		return
	}

	p.BaseHandler.SuccessResponse(c, nil, "Gán role cho user thành công")
}

// RemoveRolesFromUser - Xóa roles khỏi user
func (p *PermissionController) RemoveRolesFromUser(c *gin.Context) {
	logger.Info("controller", "RemoveRolesFromUser", "RemoveRolesFromUser controller called")
	ctx := c.Request.Context()

	userID := c.Param("id")
	if userID == "" {
		logger.Warn("controller", "RemoveRolesFromUser", "Missing user ID")
		p.BaseHandler.BadRequest(c, "user ID is required")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "RemoveRolesFromUser", fmt.Sprintf("Bind JSON error: %v", err))
		p.BaseHandler.BadRequest(c, "Xóa role khỏi user thất bại")
		return
	}

	err := p.service.RemoveRolesFromUser(ctx, userID, req)
	if err != nil {
		logger.Error("controller", "RemoveRolesFromUser", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, "Xóa role khỏi user thất bại")
		return
	}

	p.BaseHandler.SuccessResponse(c, nil, "Xóa role khỏi user thành công")
}

// GetUserPermissions - Lấy tất cả permissions của user
func (p *PermissionController) GetUserPermissions(c *gin.Context) {
	logger.Info("controller", "GetUserPermissions", "GetUserPermissions controller called")
	ctx := c.Request.Context()

	userID := c.Param("id")
	if userID == "" {
		logger.Warn("controller", "GetUserPermissions", "Missing user ID")
		p.BaseHandler.BadRequest(c, "user ID is required")
		return
	}

	permissions, err := p.service.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("controller", "GetUserPermissions", fmt.Sprintf("Service error: %v", err))
		p.BaseHandler.BadRequest(c, "Lấy permissions của user thất bại")
		return
	}

	p.BaseHandler.SuccessResponse(c, permissions, "Lấy permissions của user thành công")
}

