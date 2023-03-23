package dto

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/shop/model"
)

type SellerPromotion struct {
	model.ShopPromotion
	Status  string  `json:"status" gorm:"column:status"`
	Product Product `json:"product"`
}

func (SellerPromotion) TableName() string {
	return "shop_promotions"
}

type Product struct {
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	ImageURL string              `json:"imageUrl,omitempty"`
	SKUs     []*productModel.Sku `json:"skus,omitempty"`
}

type SellerPromotionFilterRequest struct {
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
	Status string `form:"status"`
	Name   string `form:"name"`
}

func (p *SellerPromotionFilterRequest) Validate() {
	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}
}

func (p *SellerPromotionFilterRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
