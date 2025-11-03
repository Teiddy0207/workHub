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

type RoleController struct {
	handler.BaseHandler
	service service.RoleServiceInterface
}

func NewRoleController(service service.RoleServiceInterface) *RoleController {
	return &RoleController{
		BaseHandler: handler.NewBaseHandler(),
		service:     service,
	}
}

// CreateRole - Tạo role mới
func (r *RoleController) CreateRole(c *gin.Context) {
	logger.Info("controller", "CreateRole", "CreateRole controller called")
	ctx := c.Request.Context()

	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "CreateRole", fmt.Sprintf("Bind JSON error: %v", err))
		r.BaseHandler.BadRequest(c, constant.ROLE_CREATED_FAIL)
		return
	}

	logger.Info("controller", "CreateRole", fmt.Sprintf("Request received: name=%s, code=%s", req.Name, req.Code))

	response, err := r.service.CreateRole(ctx, req)
	if err != nil {
		logger.Error("controller", "CreateRole", fmt.Sprintf("Service error: %v", err))
		// Dùng error cụ thể nếu có, ngược lại dùng FAIL constant
		if errors.Is(err, constant.ErrTakenCredential) || errors.Is(err, constant.ErrNotFound) {
			r.BaseHandler.BadRequest(c, err.Error())
		} else {
			r.BaseHandler.BadRequest(c, constant.ROLE_CREATED_FAIL)
		}
		return
	}

	logger.Info("controller", "CreateRole", "Role created successfully")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_CREATED_SUCCESSFULLY)
}

// GetRoleByID - Lấy thông tin role theo ID
func (r *RoleController) GetRoleByID(c *gin.Context) {
	logger.Info("controller", "GetRoleByID", "GetRoleByID controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "GetRoleByID", "Missing role ID")
		r.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	logger.Info("controller", "GetRoleByID", fmt.Sprintf("Getting role with ID: %s", id))

	response, err := r.service.GetRoleByID(ctx, id)
	if err != nil {
		logger.Error("controller", "GetRoleByID", fmt.Sprintf("Service error: %v", err))
		// GetRoleByID thường chỉ có lỗi ErrNotFound, giữ nguyên error message
		r.BaseHandler.BadRequest(c, err.Error())
		return
	}

	logger.Info("controller", "GetRoleByID", "Role retrieved successfully")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_GET_SUCCESSFULLY)
}

// UpdateRole - Cập nhật role
func (r *RoleController) UpdateRole(c *gin.Context) {
	logger.Info("controller", "UpdateRole", "UpdateRole controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "UpdateRole", "Missing role ID")
		r.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	var req dto.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "UpdateRole", fmt.Sprintf("Bind JSON error: %v", err))
		r.BaseHandler.BadRequest(c, constant.ROLE_UPDATED_FAIL)
		return
	}

	logger.Info("controller", "UpdateRole", fmt.Sprintf("Updating role with ID: %s", id))

	response, err := r.service.UpdateRole(ctx, id, req)
	if err != nil {
		logger.Error("controller", "UpdateRole", fmt.Sprintf("Service error: %v", err))
		// Dùng error cụ thể nếu có, ngược lại dùng FAIL constant
		if errors.Is(err, constant.ErrNotFound) || errors.Is(err, constant.ErrTakenCredential) {
			r.BaseHandler.BadRequest(c, err.Error())
		} else {
			r.BaseHandler.BadRequest(c, constant.ROLE_UPDATED_FAIL)
		}
		return
	}

	logger.Info("controller", "UpdateRole", "Role updated successfully")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_UPDATED_SUCCESSFULLY)
}

// DeleteRole - Xóa role
func (r *RoleController) DeleteRole(c *gin.Context) {
	logger.Info("controller", "DeleteRole", "DeleteRole controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "DeleteRole", "Missing role ID")
		r.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	logger.Info("controller", "DeleteRole", fmt.Sprintf("Deleting role with ID: %s", id))

	err := r.service.DeleteRole(ctx, id)
	if err != nil {
		logger.Error("controller", "DeleteRole", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) {
			r.BaseHandler.BadRequest(c, err.Error())
		} else {
			r.BaseHandler.BadRequest(c, constant.ROLE_DELETED_FAIL)
		}
		return
	}

	logger.Info("controller", "DeleteRole", "Role deleted successfully")
	r.BaseHandler.SuccessResponse(c, nil, constant.ROLE_DELETED)
}

// ListRoles - Lấy danh sách roles
func (r *RoleController) ListRoles(c *gin.Context) {
	logger.Info("controller", "ListRoles", "ListRoles controller called")
	ctx := c.Request.Context()

	params := params.NewQueryParams(c)
	logger.Info("controller", "ListRoles", fmt.Sprintf("Query params: page=%d, size=%d, search=%s", 
		params.PageNumber, params.PageSize, params.Search))

	response, err := r.service.ListRoles(ctx, params)
	if err != nil {
		logger.Error("controller", "ListRoles", fmt.Sprintf("Service error: %v", err))
		// Dùng FAIL constant cho lỗi database/internal
		r.BaseHandler.BadRequest(c, constant.ROLE_GET_LIST_FAIL)
		return
	}

	logger.Info("controller", "ListRoles", "Roles listed successfully")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_GET_LIST_SUCCESSFULLY)
}