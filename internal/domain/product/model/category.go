package model

import "gorm.io/gorm"

type Category struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	ImageURL   string      `json:"imageUrl"`
	MinPrice   *float64    `json:"minPrice,omitempty"`
	ParentID   *int        `json:"parentId,omitempty"`
	Children   []*Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	gorm.Model `json:"-"`
}
