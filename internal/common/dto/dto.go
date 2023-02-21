package dto

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
}
