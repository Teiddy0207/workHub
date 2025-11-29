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
