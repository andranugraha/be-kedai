package model

import (
	"time"

	"gorm.io/gorm"
)

type WalletHistory struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Reference string    `json:"reference"`
	Date      time.Time `json:"date" gorm:"default:CURRENT_TIMESTAMP"`
	Amount    float64   `json:"amount"`
	WalletId  int       `json:"walletId"`

	gorm.Model `json:"-"`
}

const (
	WalletHistoryTypeTopup    = "Top-up"
	WalletHistoryTypeCheckout = "Checkout"
)
