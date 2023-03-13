package dto

type ShipmentCourierResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"isActive"`
}

type MatchingProductCourierRequest struct {
	ProductIDs []int `form:"productId" binding:"required"`
	Slug       string
	ShopID     int
}

type ToggleShopCourierRequest struct {
	CourierId int `json:"courierId" binding:"required"`
}

type ToggleShopCourierResponse struct {
	CourierId int  `json:"courierId"`
	IsToggled bool `json:"isActive"`
}
