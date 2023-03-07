package model

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type InvoicePerShop struct {
	ID              int     `json:"id"`
	Code            string  `json:"code"`
	Total           float64 `json:"total"`
	Subtotal        float64 `json:"subtotal"`
	ShippingCost    float64 `json:"shippingCost"`
	VoucherAmount   float64 `json:"voucherAmount"`
	VoucherType     string  `json:"voucherType"`
	PromotionAmount float64 `json:"promotionAmount"`
	PromotionType   string  `json:"promotionType"`
	Status          string  `json:"status"`

	UserID      int `json:"userId"`
	VoucherID   int `json:"voucherId"`
	ShopID      int `json:"shopId"`
	PromotionID int `json:"promotionId"`
	InvoiceID   int `json:"invoiceId"`

	Shop       *model.Shop      `json:"shop,omitempty"`
	StatusList []*InvoiceStatus `json:"statusList"`

	gorm.Model `json:"-"`
}
