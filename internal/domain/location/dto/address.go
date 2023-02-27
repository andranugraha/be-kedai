package dto

import (
	"kedai/backend/be-kedai/internal/domain/location/model"
)

type AddressRequest struct {
	ID            int
	UserID        int
	Name          string `json:"name" binding:"required,max=30"`
	PhoneNumber   string `json:"phoneNumber" binding:"required,numeric"`
	Street        string `json:"street" binding:"required,max=200"`
	Details       string `json:"details" binding:"max=30"`
	SubdistrictID int    `json:"subdistrictId" binding:"required,numeric,min=1"`
	IsDefault     *bool  `json:"isDefault" binding:"required"`
}

func (r *AddressRequest) ValidateId() {
	if r.ID <= 0 {
		r.ID = 0
	}
}

func (r *AddressRequest) ToUserAddress() *model.UserAddress {

	if r.IsDefault == nil {
		r.IsDefault = new(bool)
		*r.IsDefault = false
	}

	return &model.UserAddress{
		ID:            r.ID,
		UserID:        r.UserID,
		Name:          r.Name,
		PhoneNumber:   r.PhoneNumber,
		Street:        r.Street,
		Details:       r.Details,
		SubdistrictID: r.SubdistrictID,
		IsDefault:     *r.IsDefault,
	}
}

func ToAddressList(addresses []*model.UserAddress, defaultAddressId *int) []*model.UserAddress {
	if len(addresses) == 0 || defaultAddressId == nil {
		return addresses
	}

	var newAddresses []*model.UserAddress
	var defaultAddressIdx int

	for i, address := range addresses {
		if address.ID == *defaultAddressId {
			address.IsDefault = true
			defaultAddressIdx = i
		}
	}

	newAddresses = append(newAddresses, addresses[defaultAddressIdx])

	for _, address := range addresses {
		if address.ID != *defaultAddressId {
			newAddresses = append(newAddresses, address)
		}
	}

	return newAddresses
}
