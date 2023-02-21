package model

import "gorm.io/gorm"

type Sku struct {
	ID        int     `json:"id"`
	Sku       string  `json:"sku"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	ProductId int     `json:"productId"`

	Product  Product   `json:"product" gorm:"foreignKey:ProductId"`
	Variants []Variant `json:"variants" gorm:"many2many:product_variants;"`

	gorm.Model `json:"-"`
}
