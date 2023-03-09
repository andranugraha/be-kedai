package dto

import (
	"kedai/backend/be-kedai/internal/domain/location/model"

	"sort"

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
	IsPickup      *bool  `json:"isPickup" binding:"required"`
}

func (r *AddressRequest) Validate() {
	if r.ID <= 0 {
		r.ID = 0
	}
}

func (r *AddressRequest) ToUserAddress() *model.UserAddress {

	if r.IsDefault == nil {
		r.IsDefault = new(bool)
		*r.IsDefault = false
	}

	if r.IsPickup == nil {
		r.IsPickup = new(bool)
		*r.IsPickup = false
	}

	return &model.UserAddress{
		ID:            r.ID,
		UserID:        r.UserID,
		Name:          r.Name,
		PhoneNumber:   r.PhoneNumber,
		Street:        r.Street,
		Details:       r.Details,
		SubdistrictID: r.SubdistrictID,
		IsDefault:     r.IsDefault,
		IsPickup:      r.IsPickup,
	}
}

func ToAddressList(addresses []*model.UserAddress, defaultAddressId *int, pickupAddressId *int) []*model.UserAddress {
	if len(addresses) == 0 || (defaultAddressId == nil && pickupAddressId == nil) {
		return addresses
	}

	trueValue := true
	falseValue := false

	for _, address := range addresses {
		if address.ID == *defaultAddressId {
			address.IsDefault = &trueValue
		} else {
			address.IsDefault = &falseValue
		}

		if address.ID == *pickupAddressId {
			address.IsPickup = &trueValue
		} else {
			address.IsPickup = &falseValue
		}

	}

	sort.Slice(addresses, func(i, j int) bool {
		if *addresses[i].IsDefault && !*addresses[j].IsDefault {
			return true
		} else if !*addresses[i].IsDefault && *addresses[j].IsDefault {
			return false
		} else {
			return *addresses[i].IsPickup && !*addresses[j].IsPickup
		}
	})

	return addresses
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
