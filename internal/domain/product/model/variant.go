package model

import (
	"gorm.io/gorm"
)

type Variant struct {
	ID       int    `json:"id"`
	Value    string `json:"value"`
	MediaUrl string `json:"mediaUrl"`
	GroupId  int    `json:"groupId"`

	Group           *VariantGroup    `json:"group,omitempty" gorm:"foreignKey:GroupId"`
	ProductVariants []ProductVariant `json:"productVariants,omitempty" gorm:"foreignKey:VariantId"`

	gorm.Model `json:"-"`
}
