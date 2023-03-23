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

type Product struct {
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	ImageURL string              `json:"imageUrl,omitempty"`
	SKUs     []*productModel.Sku `json:"skus,omitempty"`
}
