package model

import "gorm.io/gorm"

type Invoice struct {
	ID              int      `json:"id"`
	Code            string   `json:"code"`
	Subtotal        float64  `json:"subtotal"`
	Total           float64  `json:"total"`
	VoucherAmount   *float64 `json:"voucherAmount,omitempty"`
	VoucherType     *string  `json:"voucherType,omitempty"`
	UserID          int      `json:"userId"`
	VoucherID       *int     `json:"voucherId,omitempty"`
	PaymentMethodID int      `json:"courierServiceId"`

	InvoicePerShop []*InvoicePerShop `json:"invoicePerShop"`

	gorm.Model `json:"-"`
}
