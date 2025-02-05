package model

import (
	"gorm.io/gorm"
)

type UserAddress struct {
	ID            int          `json:"id"`
	UserID        int          `json:"-"`
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

	IsDefault *bool `json:"isDefault,omitempty" gorm:"<-:false"`
	IsPickup  *bool `json:"isPickup,omitempty" gorm:"<-:false"`

	gorm.Model `json:"-"`
}
