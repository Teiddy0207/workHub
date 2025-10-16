package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QueryParams - gộp tất cả vào một struct đơn giản
type QueryParams struct {
	Page     int    `json:"page" form:"page"`
	Limit    int    `json:"limit" form:"limit"`
	Keyword  string `json:"keyword" form:"keyword"`
	Total    int    `json:"total"`
	HasMore  bool   `json:"has_more"`
}

// ParseQuery - parse query params và trả về struct đơn giản
func ParseQuery(c *gin.Context) QueryParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	keyword := c.Query("keyword")

	// Validation
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return QueryParams{
		Page:    page,
		Limit:   limit,
		Keyword: keyword,
	}
}

// SetTotal - set total và tính has_more
func (q *QueryParams) SetTotal(total int) {
	q.Total = total
	q.HasMore = q.Page < int(math.Ceil(float64(total)/float64(q.Limit)))
}

// GetOffset - tính offset cho database
func (q QueryParams) GetOffset() int {
	return (q.Page - 1) * q.Limit
}

// ToMeta - chuyển thành meta cho response
func (q QueryParams) ToMeta() map[string]interface{} {
	return map[string]interface{}{
		"page":     q.Page,
		"limit":    q.Limit,
		"total":    q.Total,
		"has_more": q.HasMore,
	}
}
