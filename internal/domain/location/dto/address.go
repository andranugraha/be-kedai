package dto

import "kedai/backend/be-kedai/internal/domain/location/model"

type AddAddressRequest struct {
	UserID        int    `json:"userId"`
	Name          string `json:"name" binding:"required,max=20"`
	PhoneNumber   string `json:"phoneNumber" binding:"required,numeric"`
	Street        string `json:"street" binding:"required,max=100"`
	Details       string `json:"details" binding:"max=50"`
	SubdistrictID int    `json:"subdistrictId" binding:"required,numeric,min=1"`
	IsDefault     *bool  `json:"isDefault" binding:"required"`
}

func (r *AddAddressRequest) ToUserAddress() *model.UserAddress {
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
