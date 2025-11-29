package controller

import (
	"errors"
	"fmt"
	"workHub/constant"
	"workHub/helper"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/logger"
	"workHub/pkg/handler"
	"workHub/pkg/params"

	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	handler.BaseHandler
	service service.WorkspaceServiceInterface
}

func NewWorkspaceController(service service.WorkspaceServiceInterface) *WorkspaceController {
	return &WorkspaceController{
		BaseHandler: handler.NewBaseHandler(),
		service:     service,
	}
}

func (w *WorkspaceController) CreateWorkspace(c *gin.Context) {
	logger.Info("controller", "CreateWorkspace", "CreateWorkspace controller called")
	ctx := c.Request.Context()

	// Lấy email từ token trong header (helper function parse token)
	userEmail, err := helper.GetUserEmailFromToken(c)
	if err != nil {
		logger.Error("controller", "CreateWorkspace", fmt.Sprintf("Failed to get user email from token: %v", err))
		w.BaseHandler.BadRequest(c, "Unauthorized: Invalid token")
		return
	}

	var req dto.WorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "CreateWorkspace", fmt.Sprintf("Bind JSON error: %v", err))
		w.BaseHandler.BadRequest(c, constant.WORKSPACE_CREATED_FAIL)
		return
	}

	logger.Info("controller", "CreateWorkspace", fmt.Sprintf("Request received: name=%s, user_email=%s", req.Name, userEmail))

	// Truyền email xuống service để service tự query database và set owner_id
	response, err := w.service.CreateWorkspace(ctx, req, userEmail)
	if err != nil {
		logger.Error("controller", "CreateWorkspace", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) {
			w.BaseHandler.BadRequest(c, err.Error())
		} else {
			w.BaseHandler.BadRequest(c, constant.WORKSPACE_CREATED_FAIL)
		}
		return
	}

	logger.Info("controller", "CreateWorkspace", "Workspace created successfully")
	w.BaseHandler.SuccessResponse(c, response, constant.WORKSPACE_CREATED_SUCCESSFULLY)
}

func (w *WorkspaceController) GetWorkspaceByID(c *gin.Context) {
	logger.Info("controller", "GetWorkspaceByID", "GetWorkspaceByID controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "GetWorkspaceByID", "Missing workspace ID")
		w.BaseHandler.BadRequest(c, "workspace ID is required")
		return
	}

	logger.Info("controller", "GetWorkspaceByID", fmt.Sprintf("Getting workspace with ID: %s", id))

	response, err := w.service.GetWorkspaceByID(ctx, id)
	if err != nil {
		logger.Error("controller", "GetWorkspaceByID", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) {
			w.BaseHandler.BadRequest(c, err.Error())
		} else {
			w.BaseHandler.BadRequest(c, constant.WORKSPACE_GET_FAIL)
		}
		return
	}

	logger.Info("controller", "GetWorkspaceByID", "Workspace retrieved successfully")
	w.BaseHandler.SuccessResponse(c, response, constant.WORKSPACE_GET_SUCCESSFULLY)
}

func (w *WorkspaceController) UpdateWorkspace(c *gin.Context) {
	logger.Info("controller", "UpdateWorkspace", "UpdateWorkspace controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "UpdateWorkspace", "Missing workspace ID")
		w.BaseHandler.BadRequest(c, "workspace ID is required")
		return
	}

	var req dto.WorkspaceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "UpdateWorkspace", fmt.Sprintf("Bind JSON error: %v", err))
		w.BaseHandler.BadRequest(c, constant.WORKSPACE_UPDATED_FAIL)
		return
	}

	logger.Info("controller", "UpdateWorkspace", fmt.Sprintf("Updating workspace with ID: %s", id))

	response, err := w.service.UpdateWorkspace(ctx, id, req)
	if err != nil {
		logger.Error("controller", "UpdateWorkspace", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) {
			w.BaseHandler.BadRequest(c, err.Error())
		} else {
			w.BaseHandler.BadRequest(c, constant.WORKSPACE_UPDATED_FAIL)
		}
		return
	}

	logger.Info("controller", "UpdateWorkspace", "Workspace updated successfully")
	w.BaseHandler.SuccessResponse(c, response, constant.WORKSPACE_UPDATED_SUCCESSFULLY)
}

func (w *WorkspaceController) DeleteWorkspace(c *gin.Context) {
	logger.Info("controller", "DeleteWorkspace", "DeleteWorkspace controller called")
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		logger.Warn("controller", "DeleteWorkspace", "Missing workspace ID")
		w.BaseHandler.BadRequest(c, "workspace ID is required")
		return
	}

	logger.Info("controller", "DeleteWorkspace", fmt.Sprintf("Deleting workspace with ID: %s", id))

	err := w.service.DeleteWorkspace(ctx, id)
	if err != nil {
		logger.Error("controller", "DeleteWorkspace", fmt.Sprintf("Service error: %v", err))
		if errors.Is(err, constant.ErrNotFound) {
			w.BaseHandler.BadRequest(c, err.Error())
		} else {
			w.BaseHandler.BadRequest(c, constant.WORKSPACE_DELETED_FAIL)
		}
		return
	}

	logger.Info("controller", "DeleteWorkspace", "Workspace deleted successfully")
	w.BaseHandler.SuccessResponse(c, nil, constant.WORKSPACE_DELETED)
}

func (w *WorkspaceController) ListWorkspaces(c *gin.Context) {
	logger.Info("controller", "ListWorkspaces", "ListWorkspaces controller called")
	ctx := c.Request.Context()

	params := params.NewQueryParams(c)
	logger.Info("controller", "ListWorkspaces", fmt.Sprintf("Query params: page=%d, size=%d, search=%s",
		params.PageNumber, params.PageSize, params.Search))

	response, err := w.service.ListWorkspaces(ctx, params)
	if err != nil {
		logger.Error("controller", "ListWorkspaces", fmt.Sprintf("Service error: %v", err))
		w.BaseHandler.BadRequest(c, constant.WORKSPACE_GET_LIST_FAIL)
		return
	}

	logger.Info("controller", "ListWorkspaces", "Workspaces listed successfully")
	w.BaseHandler.SuccessResponse(c, response, constant.WORKSPACE_GET_LIST_SUCCESSFULLY)
}
