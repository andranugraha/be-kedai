package model

import "gorm.io/gorm"

type CourierService struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	CourierID int    `json:"courierId"`

	Courier *Courier `json:"courier,omitempty"`

	gorm.Model `json:"-"`
}
