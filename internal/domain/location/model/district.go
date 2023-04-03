package model

import "gorm.io/gorm"

type District struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CityID     int    `json:"cityId"`
	City       *City  `json:"city,omitempty"`
	gorm.Model `json:"-"`
}
