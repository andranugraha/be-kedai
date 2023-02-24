package model

import "gorm.io/gorm"

type VariantGroup struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ProductId int    `json:"productId"`
	
	gorm.Model `json:"-"`
}
