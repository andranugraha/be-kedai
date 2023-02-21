package model

import "gorm.io/gorm"

type VariantGroup struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	ProductID int        `json:"productId"`
	Variant   []*Variant `json:"variant" gorm:"foreignKey:GroupID"`

	gorm.Model `json:"-"`
}
