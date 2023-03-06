package dto

type FindShopRequest struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

type FindShopResponse struct {
	Slug         string  `json:"slug"`
	Name         string  `json:"name"`
	ProductCount int64   `json:"productCount"`
	Rating       float64 `json:"rating"`
	PhotoUrl     string  `json:"photoUrl"`
}

func (req *FindShopRequest) Validate() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}
}

func (req *FindShopRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}
