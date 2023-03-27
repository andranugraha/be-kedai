package model

import (
	"time"

	"gorm.io/gorm"
)

type ShopPromotion struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	StartPeriod *time.Time `json:"startPeriod,omitempty"`
	EndPeriod   *time.Time `json:"endPeriod,omitempty"`
	ShopId      int        `json:"shopId,omitempty"`

	Shop *Shop `json:"shop,omitempty" gorm:"foreignKey:ShopId"`

	gorm.Model `json:"-"`
}

const (
	PromotionTypePercent = "percent"
	PromotionTypeNominal = "nominal"
)
