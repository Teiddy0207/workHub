package controller

import (
	"fmt"
	"net/http"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/pkg/handler"
	"workHub/pkg/params"
	"workHub/constant"
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
	fmt.Printf("🎯 CreateRole controller called\n")
	ctx := c.Request.Context()

	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ Bind JSON error: %v\n", err)
		r.BaseHandler.BadRequest(c, "invalid request body")
		return
	}

	fmt.Printf("📝 Request received: name=%s, code=%s\n", req.Name, req.Code)

	response, err := r.service.CreateRole(ctx, req)
	if err != nil {
		fmt.Printf("❌ Service error: %v\n", err)
		r.BaseHandler.BadRequest(c, err.Error())
		return
	}

	fmt.Printf("✅ Role created successfully\n")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_CREATED_SUCCESSFULLY)
}

// GetRoleByID - Lấy thông tin role theo ID
func (r *RoleController) GetRoleByID(c *gin.Context) {
	fmt.Printf("🎯 GetRoleByID controller called\n")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		fmt.Printf("❌ Missing role ID\n")
		r.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	fmt.Printf("📝 Getting role with ID: %s\n", id)

	response, err := r.service.GetRoleByID(ctx, id)
	if err != nil {
		fmt.Printf("❌ Service error: %v\n", err)
		r.BaseHandler.BadRequest(c, err.Error())
		return
	}

	fmt.Printf("✅ Role retrieved successfully\n")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_GET_SUCCESSFULLY)
}

// UpdateRole - Cập nhật role
func (r *RoleController) UpdateRole(c *gin.Context) {
	fmt.Printf("🎯 UpdateRole controller called\n")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		fmt.Printf("❌ Missing role ID\n")
		r.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	var req dto.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ Bind JSON error: %v\n", err)
		r.BaseHandler.BadRequest(c, "invalid request body")
		return
	}

	fmt.Printf("📝 Updating role with ID: %s\n", id)

	response, err := r.service.UpdateRole(ctx, id, req)
	if err != nil {
		fmt.Printf("❌ Service error: %v\n", err)
		r.BaseHandler.BadRequest(c, err.Error())
		return
	}

	fmt.Printf("✅ Role updated successfully\n")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_UPDATED_SUCCESSFULLY)
}

// DeleteRole - Xóa role
func (r *RoleController) DeleteRole(c *gin.Context) {
	fmt.Printf("🎯 DeleteRole controller called\n")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		fmt.Printf("❌ Missing role ID\n")
		r.BaseHandler.BadRequest(c, "role ID is required")
		return
	}

	fmt.Printf("📝 Deleting role with ID: %s\n", id)

	err := r.service.DeleteRole(ctx, id)
	if err != nil {
		fmt.Printf("❌ Service error: %v\n", err)
		r.BaseHandler.BadRequest(c, err.Error())
		return
	}

	fmt.Printf("✅ Role deleted successfully\n")
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": constant.ROLE_DELETED,
	})
}

// ListRoles - Lấy danh sách roles
func (r *RoleController) ListRoles(c *gin.Context) {
	fmt.Printf("🎯 ListRoles controller called\n")
	ctx := c.Request.Context()

	params := params.NewQueryParams(c)
	fmt.Printf("📝 Query params: page=%d, size=%d, search=%s\n", 
		params.PageNumber, params.PageSize, params.Search)

	response, err := r.service.ListRoles(ctx, params)
	if err != nil {
		fmt.Printf("❌ Service error: %v\n", err)
		r.BaseHandler.BadRequest(c, err.Error())
		return
	}

	fmt.Printf("✅ Roles listed successfully\n")
	r.BaseHandler.SuccessResponse(c, response, constant.ROLE_GET_LIST_SUCCESSFULLY)
}