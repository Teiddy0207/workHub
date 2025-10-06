package controller

import (
    "net/http"
    "strconv"
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
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
	keyword := c.Query("keyword")

    users, meta, err := a.service.GetListUser(ctx, keyword, page, limit)
    if err != nil {
        utils.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.Success(c, "get users success", users, meta)
}