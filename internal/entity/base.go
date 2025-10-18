package entity

type GetListRequest struct {
	Limit  int    `form:"limit" json:"limit"`
	Offset int    `form:"offset" json:"offset"`
	Sort   string `form:"sort" json:"sort"` // ví dụ: "name asc", "created_at desc"
}

type Pagination[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"total_items"`
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

type FindRequestFilter struct {
	Key           string                 `json:"key"`
	Value         interface{}            `json:"value"`
	Operator      string                 `json:"operator"`
	IgnorePrepare bool                   `json:"ignore_prepare"`
	SubFilters    []FindRequestFilter    `json:"sub_filters"`
	CustomFunc    func() (string, interface{}) `json:"-"`
}