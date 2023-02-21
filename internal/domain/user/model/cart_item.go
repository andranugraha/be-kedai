package model

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type UserCartItem struct {
	ID       int     `json:"id"`
	Quantity int     `json:"quantity"`
	Notes    *string `json:"notes"`
	UserId   int     `json:"userId"`
	SkuId    int     `json:"skuId"`

	Sku productModel.Sku `json:"sku" gorm:"foreignKey:SkuId"`

	gorm.Model `json:"-"`
}
