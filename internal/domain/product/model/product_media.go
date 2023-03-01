package model

import "gorm.io/gorm"

type ProductMedia struct {
	ID        int    `json:"id"`
	Url       string `json:"url"`
	ProductID int    `json:"productId"`

	gorm.Model `json:"-"`
}

func (ProductMedia) TableName() string {
	return "product_medias"
}
