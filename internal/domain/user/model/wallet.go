package model

import "gorm.io/gorm"

type Wallet struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	Pin        string  `json:"-"`
	Balance    float64 `json:"balance"`
	gorm.Model `json:"-"`
}
