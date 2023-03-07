package model

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type Transaction struct {
	ID         int     `json:"id"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"totalPrice"`
	Quantity   int     `json:"quantity"`
	Note       *string `json:"note"`

	InvoiceID int `json:"invoiceId"`
	AddressID int `json:"addressId"`
	UserID    int `json:"userId"`
	SkuID     int `json:"skuId"`

	User *userModel.User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Sku  *productModel.Sku `json:"sku,omitempty" gorm:"foreignKey:SkuID"`

	gorm.Model `json:"-"`
}
