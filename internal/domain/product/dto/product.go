package dto

import (
	"kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
)

type ProductDetail struct {
	model.Product
	Vouchers         []*shopModel.ShopVoucher `json:"vouchers,omitempty" gorm:"->:false"`
	Couriers         []*shopModel.Courier     `json:"couriers" gorm:"->:false"`
	MinPrice         float64                  `json:"minPrice"`
	MaxPrice         float64                  `json:"maxPrice"`
	PromotionPercent *float64                 `json:"promotionPercent,omitempty"`
}

func (ProductDetail) TableName() string {
	return "products"
}

type ProductResponse struct {
	ID           int     `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	View         int     `json:"view"`
	IsHazardous  bool    `json:"isHazardous"`
	Weight       float64 `json:"weight"`
	Length       float64 `json:"length"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	PackagedSize float64 `json:"packagedSize"`
	IsNew        bool    `json:"isNew"`
	IsActive     bool    `json:"isActive"`
	Rating       float64 `json:"rating"`
	Sold         int     `json:"sold"`

	MinPrice         float64  `json:"minPrice"`
	MaxPrice         float64  `json:"maxPrice"`
	Address          string   `json:"address"`
	PromotionPercent *float64 `json:"promotionPercent,omitempty"`
	ImageURL         string   `json:"imageUrl"`
	DefaultSkuID     int      `json:"defaultSkuId"`

	ShopID     int             `json:"shopId"`
	Shop       *shopModel.Shop `json:"shop,omitempty"`
	CategoryID int             `json:"categoryId"`
}

func (ProductResponse) TableName() string {
	return "products"
}

type RecommendationByCategoryIdRequest struct {
	CategoryId int `form:"categoryId" binding:"required,gte=1"`
	ProductId  int `form:"productId" binding:"required,gte=1"`
}
