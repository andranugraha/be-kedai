package model

import (
	"time"

	"gorm.io/gorm"
)

type ShopPromotion struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	StartPeriod time.Time `json:"startPeriod"`
	EndPeriod   time.Time `json:"endPeriod"`
	ShopId      int       `json:"shopId"`

	Shop *Shop `json:"shop,omitempty" gorm:"foreignKey:ShopId"`

	gorm.Model `json:"-"`
}

const (
	PromotionTypePercent = "percent"
	PromotionTypeNominal = "nominal"
)
