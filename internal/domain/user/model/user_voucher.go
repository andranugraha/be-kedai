package model

import (
	"time"

	"gorm.io/gorm"
)

type UserVoucher struct {
	ID                   int       `json:"id"`
	IsUsed               bool      `json:"isUsed"`
	ExpiredAt            time.Time `json:"expiredAt"`
	ShopVoucherId        *int      `json:"shopVoucherId,omitempty"`
	MarketplaceVoucherId *int      `json:"marketplaceVoucherId,omitempty"`
	UserId               int       `json:"userId"`

	gorm.Model `json:"-"`
}
