package dto

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"time"
)

type SellerPromotion struct {
	model.ShopPromotion
	Status  string                        `json:"status" gorm:"column:status"`
	Product []*dto.SellerProductPromotion `json:"products"`
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

func (Product) TableName() string {
	return "products"
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

type CreateShopPromotionRequest struct {
	Name              string    `json:"name" binding:"required,min=1,max=100"`
	StartPeriod       time.Time `json:"startPeriod" binding:"required"`
	EndPeriod         time.Time `json:"endPeriod" binding:"required"`
	ProductPromotions []*dto.CreateProductPromotionRequest
}
