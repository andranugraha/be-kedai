package model

import "gorm.io/gorm"

type Wallet struct {
	ID         int     `json:"id"`
	UserID     int     `json:"userId"`
	Pin        string  `json:"-"`
	Balance    float64 `json:"balance"`
	gorm.Model `json:"-"`
}
