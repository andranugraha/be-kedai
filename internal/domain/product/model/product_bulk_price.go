package model

import "gorm.io/gorm"

type ProductBulkPrice struct {
	ID          int     `json:"id"`
	MinQuantity int     `json:"minQuantity"`
	Price       float64 `json:"price"`
	ProductID   int     `json:"productId"`

	gorm.Model `json:"-"`
}
