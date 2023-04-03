package model

import (
	"kedai/backend/be-kedai/internal/utils/random"

	"gorm.io/gorm"
)

type Wallet struct {
	ID         int     `json:"id"`
	UserID     int     `json:"userId"`
	Pin        string  `json:"-"`
	Balance    float64 `json:"balance"`
	Number     string  `json:"number"`
	gorm.Model `json:"-"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	randomGen := random.NewRandomUtils(&random.RandomUtilsConfig{})
	walletIdLength := 16
	w.Number = randomGen.GenerateNumericString(walletIdLength)
	return
}
