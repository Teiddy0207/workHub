package controller

import (
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

type BookController struct {
	handler.BaseHandler
	service service.BookServiceInterface
}

func NewBookController(service service.BookServiceInterface) *BookController {
	return &BookController{
		BaseHandler: handler.NewBaseHandler(),
		service:     service,
	}
}

func (b *BookController) Create(c *gin.Context) {
	ctx := c.Request.Context()

	// Kiểm tra quyền: chỉ admin mới được tạo sách
	userRole, err := helper.GetUserRole(c)
	if err != nil {
		userRole = "student" // Default nếu không lấy được
	}

	if userRole != "admin" {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Forbidden: Only admin can create books",
		})
		return
	}

	var req dto.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("controller", "CreateBook", fmt.Sprintf("Bind JSON error: %v", err))
		b.BaseHandler.BadRequest(c, "Tạo sách thất bại")
		return
	}

	book, err := b.service.Create(ctx, req)
	if err != nil {
		if err == constant.ErrISBNAlreadyExists {
			c.JSON(409, gin.H{
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
		"message": "Book created successfully",
		"data":    book,
	})
}

func (b *BookController) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	book, err := b.service.GetByID(ctx, id)
	if err != nil {
		if err == constant.ErrBookNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, book, "Lấy thông tin sách thành công")
}

func (b *BookController) List(c *gin.Context) {
	ctx := c.Request.Context()
	page := params.NewQueryParams(c)
	category := c.Query("category")
	search := c.Query("search")

	books, err := b.service.List(ctx, page, category, search)
	if err != nil {
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, books, "Lấy danh sách sách thành công")
}

func (b *BookController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Kiểm tra quyền: chỉ admin mới được cập nhật sách
	userRole, err := helper.GetUserRole(c)
	if err != nil {
		userRole = "student"
	}

	if userRole != "admin" {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Forbidden: Only admin can update books",
		})
		return
	}

	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.BaseHandler.BadRequest(c, "Cập nhật sách thất bại")
		return
	}

	book, err := b.service.Update(ctx, id, req)
	if err != nil {
		if err == constant.ErrBookNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err == constant.ErrISBNAlreadyExists {
			c.JSON(409, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, book, "Cập nhật sách thành công")
}

func (b *BookController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Kiểm tra quyền: chỉ admin mới được xóa sách
	userRole, err := helper.GetUserRole(c)
	if err != nil {
		userRole = "student"
	}

	if userRole != "admin" {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Forbidden: Only admin can delete books",
		})
		return
	}

	err = b.service.Delete(ctx, id)
	if err != nil {
		if err == constant.ErrCannotDeleteBookWithActiveBorrows {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Book deleted successfully",
	})
}

func (b *BookController) Search(c *gin.Context) {
	ctx := c.Request.Context()
	page := params.NewQueryParams(c)

	var searchReq dto.BookSearchRequest
	if err := c.ShouldBindQuery(&searchReq); err != nil {
		b.BaseHandler.BadRequest(c, "Tìm kiếm thất bại")
		return
	}

	books, err := b.service.Search(ctx, searchReq, page)
	if err != nil {
		b.BaseHandler.BadRequest(c, err.Error())
		return
	}

	b.BaseHandler.SuccessResponse(c, books, "Tìm kiếm sách thành công")
}

func (b *BookController) GetBorrowsByBookID(c *gin.Context) {
	// Controller này sẽ được implement trong BorrowController
	// Giữ lại để có thể dùng sau
}
