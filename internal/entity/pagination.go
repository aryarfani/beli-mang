package entity

type PaginationMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func NewPaginationMeta(limit, offset, total int) *PaginationMeta {
	return &PaginationMeta{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
