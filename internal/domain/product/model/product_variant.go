package model

import "gorm.io/gorm"

type ProductVariant struct {
	ID        int
	SkuId     int
	VariantId int

	gorm.Model
}
