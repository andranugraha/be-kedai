package model

import "gorm.io/gorm"

type Courier struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`

	Services []*CourierService `json:"services,omitempty"`

	gorm.Model `json:"-"`
}
