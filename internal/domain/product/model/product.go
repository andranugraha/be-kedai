package model

import (
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type Product struct {
	ID           int     `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	View         int     `json:"view"`
	IsHazardous  bool    `json:"isHazardous"`
	Weight       float64 `json:"weight"`
	Length       float64 `json:"length"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	PackagedSize float64 `json:"packagedSize"`
	IsNew        bool    `json:"isNew"`
	IsActive     bool    `json:"isActive"`
	Rating       float64 `json:"rating"`

	MinPrice         float64 `json:"minPrice" gorm:"<-:false"`
	MaxPrice         float64 `json:"maxPrice" gorm:"<-:false"`
	Address          string  `json:"address" gorm:"<-:false"`
	TotalSold        int     `json:"totalSold" gorm:"<-:false"`
	PromotionPercent float64 `json:"promotionPercent" gorm:"<-:false"`

	ShopID     int             `json:"shopId"`
	Shop       *shopModel.Shop `json:"shop,omitempty"`
	CategoryID int             `json:"categoryId"`

	gorm.Model `json:"-"`
}
