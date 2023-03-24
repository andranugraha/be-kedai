package model

import (
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model       `json:"-"`
	ID               int             `json:"id"`
	Message          string          `json:"message"`
	Type             string          `json:"type"`
	ShopId           int             `json:"shopId"`
	Shop             *shopModel.Shop `json:"shop" gorm:"foreignKey:ShopId"`
	UserId           int             `json:"usedId"`
	User             *userModel.User `json:"user" gorm:"foreignKey:UserId"`
	Issuer           string          `json:"issuer"`
	IsReadByOpponent bool            `json:"isReadByOpponent"`
	CreatedAt        time.Time       `json:"createdAt"`
}
