package model

import "gorm.io/gorm"

type City struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ProvinceID int       `json:"province_id"`
	Province   *Province `json:"province,omitempty"`
	gorm.Model `json:"-"`
}
