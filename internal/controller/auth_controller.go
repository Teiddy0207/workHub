package controller

import (
	"net/http"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (a *AuthController) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := a.service.Register(ctx, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, "registered success", res, nil)
}

func (a *AuthController) GetListUser(c *gin.Context) {
    ctx := c.Request.Context()
    query := utils.ParseQueryParams(c)

    items, meta, err := a.service.GetListUser(ctx, query)
    if err != nil {
        utils.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.Success(c, "Get users success", items, meta)
}