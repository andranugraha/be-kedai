package model

import "gorm.io/gorm"

type ShopCategory struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ShopId   int    `json:"shopId"`
	IsActive bool   `json:"isActive" gorm:"default:true"`

	Products []*ShopCategoryProduct `json:"products,omitempty" gorm:"foreignKey:ShopCategoryId"`

	gorm.Model `json:"-"`
}

type ShopCategoryProduct struct {
	ID             int `json:"id"`
	ShopCategoryId int `json:"shopCategoryId"`
	ProductId      int `json:"productId"`

	gorm.Model `json:"-"`
}
