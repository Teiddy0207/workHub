package utils

import "math"

type Pagination struct {
    Page       int `json:"page"`
    Limit      int `json:"limit"`
    TotalRows  int `json:"total_rows"`
    TotalPages int `json:"total_pages"`
}

func Paginate(page, limit, total int) Pagination {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    return Pagination{
        Page:       page,
        Limit:      limit,
        TotalRows:  total,
        TotalPages: totalPages,
    }
}
