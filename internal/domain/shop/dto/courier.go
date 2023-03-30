package dto

type ShipmentCourierFilterRequest struct {
	Status string `form:"status"`
}

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
	IsToggled bool `json:"isToggled"`
}

type ShipmentCourierRequest struct {
	Name    string                   `json:"name" binding:"required"`
	Code    string                   `json:"code" binding:"required"`
	Service []ShipmentCourierService `json:"service" binding:"required"`
}

type ShipmentCourierService struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	MinDuration int    `json:"minDuration"`
	MaxDuration int    `json:"maxDuration"`
}
