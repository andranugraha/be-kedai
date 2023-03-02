package model

import "gorm.io/gorm"

type Variant struct {
	ID       int    `json:"id"`
	Value    string `json:"value"`
	MediaUrl string `json:"mediaUrl"`
	GroupId  int    `json:"groupId"`

	Group *VariantGroup `json:"group,omitempty" gorm:"foreignKey:GroupId"`

	gorm.Model `json:"-"`
}
