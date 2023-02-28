package model

import (
	"time"

	"gorm.io/gorm"
)

type SealabsPay struct {
	ID         int       `json:"id"`
	CardNumber string    `json:"cardNumber"`
	CardName   string    `json:"cardName"`
	ExpiryDate time.Time `json:"expiryDate"`

	UserID     int `json:"-"`
	gorm.Model `json:"-"`
}
