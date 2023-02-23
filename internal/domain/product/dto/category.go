package dto

type GetCategoriesRequest struct {
	Depth     int  `form:"depth"`
	ParentID  int  `form:"parentId"`
	WithPrice bool `form:"withPrice"`
	Limit     int  `form:"limit"`
	Page      int  `form:"page"`
}

func (r *GetCategoriesRequest) Validate() {
	if r.Limit < 0 {
		r.Limit = 0
	}
	if r.Page < 1 {
		r.Page = 1
	}
}

func (r *GetCategoriesRequest) Offset() int {
	return int((r.Page - 1) * r.Limit)
}
