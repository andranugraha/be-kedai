package model

import (
	"time"

	"gorm.io/gorm"
)

type WalletHistory struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Reference string    `json:"reference"`
	Date      time.Time `json:"date"`
	Amount    float64   `json:"amount"`
	WalletId  int       `json:"walletId"`

	gorm.Model `json:"-"`
}
