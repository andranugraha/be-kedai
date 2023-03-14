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

type ShopFinanceOverviewResponse struct {
	ToRelease float64             `json:"toRelease"`
	Released  ShopFinanceReleased `json:"released"`
}

type ShopFinanceReleased struct {
	Week  float64 `json:"week"`
	Month float64 `json:"month"`
	Total float64 `json:"total"`
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

type CreateShopRequest struct {
	Name       string `json:"string" binding:"required,min=5,max=30"`
	AddressID  int    `json:"addressId" binding:"omitempty,gte=1"`
	CourierIDs []int  `json:"courierIds" binding:"required,dive,gte=1"`
}
