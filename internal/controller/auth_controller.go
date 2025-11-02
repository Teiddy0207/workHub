package controller

import (
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/pkg/handler"
	"workHub/pkg/params"
	"workHub/logger"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	handler.BaseHandler
	service service.AuthServiceInterface
}

func NewAuthController(service service.AuthServiceInterface) *AuthController {
	return &AuthController{
		BaseHandler: handler.NewBaseHandler(),
		service:     service}
}



func (a *AuthController) GetListUser(c *gin.Context) {
	ctx := c.Request.Context()

	params := params.NewQueryParams(c)

	users, err := a.service.GetListUser(ctx, params)

	if err != nil {
		a.BaseHandler.BadRequest(c, "get all user failed")
		return	
	}

	a.BaseHandler.SuccessResponse(c, users,  "get all user success")
}

func (a *AuthController) Login(c *gin.Context) {
	logger.Info("controller", "Login", "Login controller called")
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "Login", fmt.Sprintf("Bind JSON error: %v", err))
		a.BaseHandler.BadRequest(c, "invalid request body")
		return
	}
	
	logger.Info("controller", "Login", fmt.Sprintf("Request received: email=%s", req.Email))

	response, err := a.service.Login(ctx, req)
	if err != nil {
		logger.Error("controller", "Login", fmt.Sprintf("Service error: %v", err))
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	logger.Info("controller", "Login", "Service success, sending response")
	a.BaseHandler.SuccessResponse(c, response, "login success")
}
