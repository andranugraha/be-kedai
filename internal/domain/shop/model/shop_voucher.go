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
	IsHidden     bool      `json:"isHidden"`
	Description  string    `json:"description"`
	MinimumSpend float64   `json:"minimumSpend"`
	ExpiredAt    time.Time `json:"expiredAt"`
	ShopId       int       `json:"shopId"`

	gorm.Model `json:"-"`
}
