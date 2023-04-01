package dto

type GetSellerCategoriesRequest struct {
	Status string `form:"status"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

func (r *GetSellerCategoriesRequest) Validate() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Limit < 1 || r.Limit > 10 {
		r.Limit = 10
	}

	if r.Status != "enabled" && r.Status != "disabled" {
		r.Status = ""
	}
}

func (r *GetSellerCategoriesRequest) Offset() int {
	return (r.Page - 1) * r.Limit
}

type ShopCategory struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	ShopId       int                    `json:"shopId"`
	IsActive     bool                   `json:"isActive"`
	TotalProduct int                    `json:"totalProduct"`
	Products     []*ShopCategoryProduct `json:"products,omitempty" gorm:"-"`
}

type ShopCategoryProduct struct {
	ID       int     `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	MinPrice float64 `json:"minPrice"`
	MaxPrice float64 `json:"maxPrice"`
	ImageUrl string  `json:"imageUrl"`
	Stock    int     `json:"stock"`
}
