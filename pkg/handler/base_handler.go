package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	SuccessResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	ErrorResponse struct {
		Status  string `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details,omitempty"`
	}

	ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	ValidationResponse struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Errors  []ValidationError `json:"errors"`
	}
)

// Response handler interface and implementation
type BaseHandler interface {
	SuccessResponse(c *gin.Context, data any, message string)
	ErrorResponse(c *gin.Context, code int, message string, details any)
	BadRequest(c *gin.Context, message string, details ...any) *gin.Error
}

type responseHandler struct{}

func NewBaseHandler() BaseHandler {
	return &responseHandler{}
}

// Success response functions
func NewSuccessResponse(data any, message string) *SuccessResponse {
	return &SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// Error response functions
func NewErrorResponse(code int, message string, details ...any) *ErrorResponse {
	err := &ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}

	return err
}

func (h *responseHandler) SuccessResponse(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, NewSuccessResponse(data, message))
}

func (h *responseHandler) ErrorResponse(c *gin.Context, code int, message string, details any) {
	c.JSON(code, NewErrorResponse(code, message, details))
}

func (h *responseHandler) BadRequest(c *gin.Context, message string, details ...any) *gin.Error {
	c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, message, details))
	return &gin.Error{
		Type: gin.ErrorTypePublic,
		Err:  errors.New(message),
		Meta: details,
	}
}
