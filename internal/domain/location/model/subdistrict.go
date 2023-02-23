package model

import "gorm.io/gorm"

type Subdistrict struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PostalCode string `json:"postalCode"`
	DistrictID int    `json:"districtId"`
	gorm.Model `json:"-"`
}
