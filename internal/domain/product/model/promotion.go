package model

import (
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type ProductPromotion struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Stock         int     `json:"stock"`
	IsActive      bool    `json:"isActive"`
	PurchaseLimit int     `json:"purchaseLimit"`
	SkuId         int     `json:"skuId"`
	PromotionId   int     `json:"promotionId"`

	ShopPromotion shopModel.ShopPromotion `json:"shopPromotion" gorm:"foreignKey:PromotionId"`

	gorm.Model `json:"-"`
}
