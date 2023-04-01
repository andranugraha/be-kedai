package model

import "gorm.io/gorm"

type ShopCategory struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ShopId   int    `json:"shopId"`
	IsActive bool   `json:"isActive"`

	gorm.Model `json:"-"`
}
