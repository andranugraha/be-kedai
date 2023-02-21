package dto

type GetCitiesRequest struct {
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
	ProvinceID int    `form:"provinceId"`
	Sort       string `form:"sort"`
}

func (r *GetCitiesRequest) Validate() {
	if r.Limit < 0 {
		r.Limit = 0
	}
	if r.Page < 1 {
		r.Page = 1
	}
}

func (r *GetCitiesRequest) Offset() int {
	return int((r.Page - 1) * r.Limit)
}
