package controller

import (
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/pkg/handler"
	"workHub/pkg/params"
	"workHub/logger"
	"workHub/helper"
	"github.com/gin-gonic/gin"
	"workHub/constant"
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
		a.BaseHandler.BadRequest(c, constant.LOGIN_SUCCESSFULLY)
		return
	}
	
	logger.Info("controller", "Login", fmt.Sprintf("Request received: email=%s", req.Email))

	// Lấy IP address và User-Agent từ request
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := a.service.Login(ctx, req, ipAddress, userAgent)
	if err != nil {
		logger.Error("controller", "Login", fmt.Sprintf("Service error: %v", err))
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	logger.Info("controller", "Login", "Service success, sending response")
	a.BaseHandler.SuccessResponse(c, response, "login success")
}

func (a *AuthController) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "Register", fmt.Sprintf("Bind JSON error: %v", err))
		a.BaseHandler.BadRequest(c, "Đăng ký thất bại")
		return
	}

	response, err := a.service.Register(ctx, req)
	if err != nil {
		logger.Error("controller", "Register", fmt.Sprintf("Service error: %v", err))
		if err == constant.ErrEmailAlreadyExists {
			c.JSON(409, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data":    response,
	})
}

func (a *AuthController) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	currentUserID, _ := helper.GetUserID(c)
	currentUserRole, _ := helper.GetUserRole(c)

	// Kiểm tra quyền: admin hoặc chính user đó
	if currentUserRole != "admin" && currentUserID != userID {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		return
	}

	user, err := a.service.GetUserByID(ctx, userID)
	if err != nil {
		if err == constant.ErrUserNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	a.BaseHandler.SuccessResponse(c, user, "Lấy thông tin user thành công")
}

func (a *AuthController) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	currentUserID, _ := helper.GetUserID(c)
	currentUserRole, _ := helper.GetUserRole(c)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.BaseHandler.BadRequest(c, "Cập nhật user thất bại")
		return
	}

	user, err := a.service.UpdateUser(ctx, userID, req, currentUserID, currentUserRole)
	if err != nil {
		if err == constant.ErrForbidden {
			c.JSON(403, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrUserNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	a.BaseHandler.SuccessResponse(c, user, "Cập nhật user thành công")
}

func (a *AuthController) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	err := a.service.DeleteUser(ctx, userID)
	if err != nil {
		if err == constant.ErrCannotDeleteUserWithActiveBorrows {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

func (a *AuthController) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := helper.GetUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	user, err := a.service.GetMe(ctx, userID)
	if err != nil {
		a.BaseHandler.BadRequest(c, err.Error())
		return
	}

	a.BaseHandler.SuccessResponse(c, user, "Lấy thông tin user thành công")
}