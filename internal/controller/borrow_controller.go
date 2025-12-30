package controller

import (
	"fmt"
	"workHub/internal/dto"
	"workHub/internal/service"
	"workHub/pkg/handler"
	"workHub/pkg/params"
	"workHub/logger"
	"workHub/helper"
	"workHub/constant"

	"github.com/gin-gonic/gin"
)

type BorrowController struct {
	handler.BaseHandler
	service service.BorrowServiceInterface
}

func NewBorrowController(service service.BorrowServiceInterface) *BorrowController {
	return &BorrowController{
		BaseHandler: handler.NewBaseHandler(),
		service:     service,
	}
}

func (b *BorrowController) BorrowBook(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := helper.GetUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	var req dto.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "BorrowBook", fmt.Sprintf("Bind JSON error: %v", err))
		b.BaseHandler.BadRequest(c, "Mượn sách thất bại")
		return
	}

	borrow, err := b.service.BorrowBook(ctx, userID, req)
	if err != nil {
		if err == constant.ErrNoAvailableCopies {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrUserHasMaxBorrows {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrBookAlreadyBorrowed {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrInvalidDaysToBorrow {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Book borrowed successfully",
		"data":    borrow,
	})
}

func (b *BorrowController) ReturnBook(c *gin.Context) {
	ctx := c.Request.Context()
	borrowID := c.Param("id")

	userID, _ := helper.GetUserID(c)
	userRole, _ := helper.GetUserRole(c)

	var req dto.ReturnBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Notes là optional, không bắt buộc
	}

	borrow, err := b.service.ReturnBook(ctx, borrowID, userID, userRole, req)
	if err != nil {
		if err == constant.ErrForbidden {
			c.JSON(403, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrBorrowNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrBorrowAlreadyReturned {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, borrow, "Trả sách thành công")
}

func (b *BorrowController) GetByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	targetUserID := c.Param("id")
	status := c.Query("status")
	page := params.NewQueryParams(c)

	userID, _ := helper.GetUserID(c)
	userRole, _ := helper.GetUserRole(c)

	borrows, err := b.service.GetByUserID(ctx, userID, targetUserID, userRole, status, page)
	if err != nil {
		if err == constant.ErrForbidden {
			c.JSON(403, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, borrows, "Lấy lịch sử mượn thành công")
}

func (b *BorrowController) GetByBookID(c *gin.Context) {
	ctx := c.Request.Context()
	bookID := c.Param("id")
	status := c.Query("status")
	page := params.NewQueryParams(c)

	borrows, err := b.service.GetByBookID(ctx, bookID, status, page)
	if err != nil {
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, borrows, "Lấy lịch sử mượn thành công")
}

func (b *BorrowController) GetActiveBorrows(c *gin.Context) {
	ctx := c.Request.Context()
	status := c.Query("status")
	page := params.NewQueryParams(c)

	userID, _ := helper.GetUserID(c)
	userRole, _ := helper.GetUserRole(c)

	borrows, err := b.service.GetActiveBorrows(ctx, userID, userRole, status, page)
	if err != nil {
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, borrows, "Lấy danh sách sách đang mượn thành công")
}

