package model

import "gorm.io/gorm"

type Product struct {
	ID           int               `json:"id"`
	Code         string            `json:"code"`
	Name         string            `json:"email"`
	Description  string            `json:"description"`
	View         int               `json:"view"`
	IsHazardous  bool              `json:"isHazardous"`
	Weight       float64           `json:"weight"`
	Length       float64           `json:"length"`
	Width        float64           `json:"width"`
	Height       float64           `json:"height"`
	PackagedSize float64           `json:"packagedSize"`
	IsNew        bool              `json:"isNew"`
	IsActive     bool              `json:"isActive"`
	Rating       float64           `json:"rating"`
	ShopID       int               `json:"shopId"`
	CategoryId   int               `json:"categoryId"`
	BulkPrice    *ProductBulkPrice `json:"bulkPrice"`
	VariantGroup []*VariantGroup   `json:"variantGroup"`
	Media        []*ProductMedia   `json:"media"`

	gorm.Model `json:"-"`
}
