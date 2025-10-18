package controller

import (
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/pkg/handler"
	"workHub/pkg/params"
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
	fmt.Printf("🎯 Login controller called\n")
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ Bind JSON error: %v\n", err)
		a.BaseHandler.BadRequest(c, "invalid request body")
		return
	}
	
	fmt.Printf("📝 Request received: email=%s, password=%s\n", req.Email, req.Password)

	response, err := a.service.Login(ctx, req)
	if err != nil {
		fmt.Printf("❌ Service error: %v\n", err)
		a.BaseHandler.BadRequest(c, "login failed")
		return
	}

	fmt.Printf("✅ Service success, sending response\n")
	fmt.Printf("📤 Response data: %+v\n", response)
	a.BaseHandler.SuccessResponse(c, response, "login success")
}
