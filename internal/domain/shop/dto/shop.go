package dto

type FindShopRequest struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

func (req *FindShopRequest) Validate() {
	if req.Page < 0 {
		req.Page = 1
	}

	if req.Limit < 0 {
		req.Limit = 10
	}
}
