package params

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryParams struct {
	PageNumber     int               `form:"pageNumber"`
	PageSize       int               `form:"pageSize"`
	Search         string            `form:"search"`
	Phone          string            `form:"phone"`
	Filter         map[string]string `form:"filter"`
}

func NewQueryParams(c *gin.Context) QueryParams {
	var params QueryParams

	pageNumberStr := c.Query("pageNumber")
	pageSizeStr := c.Query("pageSize")

	if pageNumberStr != "" {
		if pageNumber, err := strconv.Atoi(pageNumberStr); err == nil && pageNumber > 0 {
			params.PageNumber = pageNumber
		} else {
			params.PageNumber = 1
		}
	} else {
		params.PageNumber = 1
	}

	if pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params.PageSize = pageSize
		} else {
			params.PageSize = 10
		}
	} else {
		params.PageSize = 10
	}

	_ = c.ShouldBindQuery(&params)

	if params.Filter == nil {
		params.Filter = make(map[string]string)
	}

	return params
}

func parseIDs(idsStr string) []int {
	var ids []int
	parts := strings.Split(idsStr, ",")
	for _, part := range parts {
		if id, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
