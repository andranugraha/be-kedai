package model

import "gorm.io/gorm"

type CourierService struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	CourierID   int    `json:"courierId"`
	MinDuration int    `json:"minDuration"`
	MaxDuration int    `json:"maxDuration"`

	Courier *Courier `json:"courier,omitempty"`

	gorm.Model `json:"-"`
}
