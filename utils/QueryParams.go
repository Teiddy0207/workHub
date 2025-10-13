package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueryParams struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Keyword string `json:"keyword"`
}

func ParseQueryParams(c *gin.Context) QueryParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	keyword := c.Query("keyword")

	return QueryParams{
		Page:    page,
		Limit:   limit,
		Keyword: keyword,
	}
}
