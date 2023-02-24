package model

import "gorm.io/gorm"

type Province struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	gorm.Model `json:"-"`
}
