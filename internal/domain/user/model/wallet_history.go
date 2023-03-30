package model

import (
	"kedai/backend/be-kedai/internal/utils/random"
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
	WalletHistoryTypeTopup      = "Top-up"
	WalletHistoryTypeCheckout   = "Checkout"
	WalletHistoryTypeWithdrawal = "Withdrawal"
	WalletHistoryTypeRefund     = "Refund"
)

func (wh *WalletHistory) BeforeCreate(tx *gorm.DB) (err error) {

	wh.Date = time.Now()
	if wh.Reference == "" {
		r := random.NewRandomUtils(&random.RandomUtilsConfig{})
		wh.Reference = r.GenerateNumericString(5)
	}
	return
}
