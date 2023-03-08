package model

import (
	locationModel "kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type InvoicePerShop struct {
	ID            int     `json:"id"`
	Code          string  `json:"code"`
	Total         float64 `json:"total"`
	Subtotal      float64 `json:"subtotal"`
	ShippingCost  float64 `json:"shippingCost"`
	VoucherAmount float64 `json:"voucherAmount"`
	VoucherType   string  `json:"voucherType"`
	Status        string  `json:"status"`

	UserID           int `json:"userId"`
	VoucherID        int `json:"voucherId"`
	AddressID        int `json:"addressId"`
	ShopID           int `json:"shopId"`
	InvoiceID        int `json:"invoiceId"`
	CourierServiceID int `json:"courierServiceId"`

	Shop           *model.Shop                `json:"shop,omitempty"`
	StatusList     []*InvoiceStatus           `json:"statusList,omitempty"`
	Address        *locationModel.UserAddress `json:"address,omitempty"`
	CourierService *model.CourierService      `json:"courierService,omitempty"`

	gorm.Model `json:"-"`
}
