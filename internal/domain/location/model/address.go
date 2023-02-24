package model

import (
	"gorm.io/gorm"
)

type UserAddress struct {
	ID            int          `json:"id"`
	UserID        int          `json:"userId"`
	Name          string       `json:"name"`
	PhoneNumber   string       `json:"phoneNumber"`
	Street        string       `json:"street"`
	Details       string       `json:"details"`
	SubdistrictID int          `json:"subdistrictId"`
	Subdistrict   *Subdistrict `json:"subdistrict,omitempty"`
	DistrictID    int          `json:"districtId"`
	District      *District    `json:"district,omitempty"`
	CityID        int          `json:"cityId"`
	City          *City        `json:"city,omitempty"`
	ProvinceID    int          `json:"provinceId"`
	Province      *Province    `json:"province,omitempty"`
	gorm.Model    `json:"-"`
}
