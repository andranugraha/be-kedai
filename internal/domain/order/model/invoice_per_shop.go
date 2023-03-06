package model

import "gorm.io/gorm"

type InvoicePerShop struct {
	ID              int      `json:"id"`
	Code            string   `json:"code"`
	Total           float64  `json:"total"`
	Subtotal        float64  `json:"subtotal"`
	ShippingCost    float64  `json:"shippingCost"`
	VoucherAmount   *float64 `json:"voucherAmount,omitempty"`
	VoucherType     *string  `json:"voucherType,omitempty"`
	PromotionAmount *float64 `json:"promotionAmount,omitempty"`
	PromotionType   *string  `json:"promotionType,omitempty"`
	Status          string   `json:"status"`

	UserID      int  `json:"userId"`
	VoucherID   *int `json:"voucherId,omitempty"`
	ShopID      int  `json:"shopId"`
	PromotionID *int `json:"promotionId,omitempty"`
	InvoiceID   int  `json:"invoiceId"`

	Transactions []Transaction `json:"transactions" gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	gorm.Model `json:"-"`
}
