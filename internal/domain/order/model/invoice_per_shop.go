package model

import "gorm.io/gorm"

type InvoicePerShop struct {
	ID               int      `json:"id"`
	Code             string   `json:"code"`
	Subtotal         float64  `json:"subtotal"`
	Total            float64  `json:"total"`
	ShippingCost     float64  `json:"shippingCost"`
	TrackingNumber   int      `json:"trackingNumber"`
	PromotionAmount  *float64 `json:"promotionAmount,omitempty"`
	PromotionType    *string  `json:"promotionType,omitempty"`
	VoucherAmount    *float64 `json:"voucherAmount,omitempty"`
	VoucherType      *string  `json:"voucherType,omitempty"`
	Status           string   `json:"status"`
	PromotionID      *int     `json:"promotionId,omitempty"`
	ShopID           int      `json:"shopId"`
	UserID           int      `json:"userId"`
	VoucherID        *int     `json:"voucherId,omitempty"`
	InvoiceID        int      `json:"invoiceId"`
	CourierServiceID int      `json:"courierServiceId"`

	gorm.Model `json:"-"`
}
