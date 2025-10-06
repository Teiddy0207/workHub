package controller

import (
	"net/http"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc service.AuthService
}

func NewAuthController(svc service.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

func (a *AuthController) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := a.svc.Register(ctx, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, "registered", res, nil)
}
