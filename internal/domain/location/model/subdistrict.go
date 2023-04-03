package model

import "gorm.io/gorm"

type Subdistrict struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	PostalCode string    `json:"postalCode"`
	DistrictID int       `json:"districtId"`
	District   *District `json:"district,omitempty"`
	gorm.Model `json:"-"`
}
