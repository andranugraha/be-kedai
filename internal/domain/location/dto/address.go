package dto

import (
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type AddressRequest struct {
	ID            int
	UserID        int
	Name          string `json:"name" binding:"required,max=30"`
	PhoneNumber   string `json:"phoneNumber" binding:"required,numeric,min=10,max=15"`
	Street        string `json:"street" binding:"required,max=200"`
	Details       string `json:"details" binding:"max=30"`
	SubdistrictID int    `json:"subdistrictId" binding:"required,numeric,min=1"`
	IsDefault     *bool  `json:"isDefault" binding:"required"`
}

type SearchAddressRequest struct {
	Keyword string `form:"keyword" binding:"required"`
}

type SearchAddressResponse struct {
	PlaceID     string `json:"placeId"`
	Description string `json:"description"`
}

type GetAddressDetailResponse struct {
	ID          int                `json:"id"`
	PlaceID     string             `json:"placeId"`
	Latitude    float64            `json:"latitude"`
	Longitude   float64            `json:"longitude"`
	Street      string             `json:"street"`
	PostalCode  string             `json:"postalCode"`
	Subdistrict *model.Subdistrict `json:"subdistrict,omitempty"`
	District    *model.District    `json:"district,omitempty"`
	City        *model.City        `json:"city,omitempty"`
	Province    *model.Province    `json:"province,omitempty"`
	gorm.Model  `json:"-"`
}
