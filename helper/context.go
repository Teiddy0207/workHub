package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"workHub/constant"
	"workHub/pkg/logger"
	"workHub/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContextGin struct {
	*gin.Context
}

type dataResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

type ListResponse struct {
	Items  interface{} `json:"data,omitempty"`
	Total  int         `json:"total"`
	Params interface{} `json:"params,omitempty"`
}

type HandlerFunc func(ctx *ContextGin)

func WithContext(hander HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hander(&ContextGin{
			ctx,
		})
	}
}

func (ctx *ContextGin) WithBody(req interface{}) bool {
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("WithBody err", err)
		ctx.JSON(http.StatusBadRequest, ctx.FormatError(err))
		return false
	}
	return true
}

func (ctx *ContextGin) FormatError(err error) gin.H {
	resp := gin.H{
		"data":    nil,
		"status":  "ERROR",
		"message": err.Error(),
	}

	var code string

	if errors.Is(err, io.EOF) {
		code = "requestInvalid"
	}

	if _, ok := err.(*json.UnmarshalTypeError); ok {
		code = "requestInvalid"
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, constant.ErrNotFound) {
		code = "NotFound"
	}

	if code != "" {
		resp["errorCode"] = code
	}

	return resp
}

func (ctx *ContextGin) OKResponse(data interface{}, message string) {
	ctx.JSON(http.StatusOK, &dataResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func (ctx *ContextGin) BadLogic(err error, msg string) {
	statusCode := http.StatusUnprocessableEntity

	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, constant.ErrNotFound) {
		statusCode = http.StatusNotFound
	}

	if errors.Is(err, constant.ErrInternalServer) {
		statusCode = http.StatusInternalServerError
	}
	resp := ctx.FormatError(err)
	if len(msg) > 0 {
		resp["message"] = msg[0]
	}

	ctx.JSON(statusCode, &dataResponse{
		Status:  "error",
		Message: msg,
	})
}

func (ctx *ContextGin) InternalServerError(err error) {
	ctx.JSON(http.StatusInternalServerError, ctx.FormatError(err))
}

func (ctx *ContextGin) OKListResponse(items interface{}, total int, params ...interface{}) {
	if items == nil {
		items = []interface{}{}
	}

	if valueType := reflect.TypeOf(items); valueType.Kind() == reflect.Slice {
		if reflect.ValueOf(items).Len() == 0 {
			items = []interface{}{}
		}
	}

	ctx.JSON(http.StatusOK, &ListResponse{Items: items, Total: total, Params: params})
}

func (ctx *ContextGin) GetPagingInfo() (int64, int64) {
	if isAll := ctx.ParseBoolQuery("isAll"); isAll != nil && *isAll {
		return 0, -1
	}

	offset, _ := strconv.Atoi(ctx.TrimDefaultQuery("offset", "0"))
	page, _ := strconv.Atoi(ctx.TrimDefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.TrimDefaultQuery("limit", "10"))
	if limit > 100 {
		limit = 10
	}

	if offset == 0 {
		offset = (page - 1) * limit
	}

	return int64(offset), int64(limit)
}

func (ctx *ContextGin) ParseBoolQuery(key string) *bool {
	val := strings.ToLower(strings.Trim(ctx.DefaultQuery(key, ""), " "))
	if !utils.Includes[string]([]string{constant.TRUE, constant.FALSE}, val) {
		return nil
	}
	result := val == constant.TRUE
	return &result
}

func (ctx *ContextGin) TrimDefaultQuery(key, defaultValue string) string {
	return strings.Trim(ctx.DefaultQuery(key, defaultValue), " ")
}
