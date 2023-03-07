package model

import "gorm.io/gorm"

type Invoice struct {
	ID            int     `json:"id"`
	Code          string  `json:"code"`
	Total         float64 `json:"total"`
	Subtotal      float64 `json:"subtotal"`
	VoucherAmount float64 `json:"voucherAmount"`
	VoucherType   string  `json:"voucherType"`

	UserID          int `json:"userId"`
	VoucherID       int `json:"voucherId"`
	PaymentMethodID int `json:"paymentMethodId"`

	InvoicePerShops []InvoicePerShop `json:"invoicePerShops" gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	gorm.Model `json:"-"`
}
