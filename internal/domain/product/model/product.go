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
	Sold         int     `json:"sold"`

	// MinPrice         float64  `json:"minPrice" gorm:"<-:false"`
	// MaxPrice         float64  `json:"maxPrice" gorm:"<-:false"`
	// Address          string   `json:"address" gorm:"<-:false"`
	// PromotionPercent *float64 `json:"promotionPercent,omitempty" gorm:"<-:false"`
	// ImageURL         string   `json:"imageUrl" gorm:"<-:false"`

	ShopID       int               `json:"shopId"`
	Shop         *shopModel.Shop   `json:"shop,omitempty"`
	CategoryID   int               `json:"categoryId"`
	Bulk         *ProductBulkPrice `json:"bulk,omitempty"`
	VariantGroup []*VariantGroup   `json:"variantGroups"`
	Media        []*ProductMedia   `json:"media"`

	gorm.Model `json:"-"`
}
