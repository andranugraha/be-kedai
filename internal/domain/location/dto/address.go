package dto

import "kedai/backend/be-kedai/internal/domain/location/model"

type AddAddressRequest struct {
	UserID        int
	Name          string `json:"name" binding:"required,max=30"`
	PhoneNumber   string `json:"phoneNumber" binding:"required,numeric"`
	Street        string `json:"street" binding:"required,max=200"`
	Details       string `json:"details" binding:"max=30"`
	SubdistrictID int    `json:"subdistrictId" binding:"required,numeric,min=1"`
	IsDefault     *bool  `json:"isDefault" binding:"required"`
}

func (r *AddAddressRequest) ToUserAddress() *model.UserAddress {

	if r.IsDefault == nil {
		r.IsDefault = new(bool)
		*r.IsDefault = false
	}

	return &model.UserAddress{
		UserID:        r.UserID,
		Name:          r.Name,
		PhoneNumber:   r.PhoneNumber,
		Street:        r.Street,
		Details:       r.Details,
		SubdistrictID: r.SubdistrictID,
		IsDefault:     *r.IsDefault,
	}
}
