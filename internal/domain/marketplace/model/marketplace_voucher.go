package model

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceVoucher struct {
	ID           int       `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type"`
	IsHidden     bool      `json:"isHidden"`
	Description  string    `json:"description"`
	MinimumSpend float64   `json:"minimumSpend"`
	ExpiredAt    time.Time `json:"expiredAt"`

	CategoryID      *int `json:"categoryId"`
	PaymentMethodID *int `json:"paymentMethodId"`

	gorm.Model `json:"-"`
}

const (
	VoucherTypePercent  = "percent"
	VoucherTypeNominal  = "nominal"
	VoucherTypeShipping = "shipping"
)
