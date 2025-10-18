package controller

import (
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
