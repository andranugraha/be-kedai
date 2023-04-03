package model

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type CartItem struct {
	ID       int    `json:"id"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
	UserId   int    `json:"-"`
	SkuId    int    `json:"skuId"`

	Sku productModel.Sku `json:"sku" gorm:"foreignKey:SkuId"`

	gorm.Model `json:"-"`
}
