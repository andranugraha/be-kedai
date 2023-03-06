package model

import (
	"time"

	"gorm.io/gorm"
)

type ShopVoucher struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type"`
	IsHidden     bool      `json:"is_hidden"`
	Description  string    `json:"description"`
	MinimumSpend float64   `json:"minimum_spend"`
	ExpiredAt    time.Time `json:"expired_at"`
	ShopId       int       `json:"shop_id"`

	gorm.Model `json:"-"`
}

const (
	VoucherTypePercent = "percent"
	VoucherTypeNominal = "nominal"
)
