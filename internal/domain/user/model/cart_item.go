package model

import (
	"gorm.io/gorm"
)

type UserCartItem struct {
	ID       int     `json:"id"`
	Quantity int     `json:"quantity"`
	Notes    *string `json:"notes"`
	UserId   int     `json:"userId"`
	SkuId    int     `json:"skuId"`

	gorm.Model `json:"-"`
}
