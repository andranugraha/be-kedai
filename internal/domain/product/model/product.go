package model

import "gorm.io/gorm"

type Product struct {
	ID           int               `json:"id"`
	Code         string            `json:"code"`
	Name         string            `json:"email"`
	Description  string            `json:"description"`
	View         int               `json:"view"`
	IsHazardous  bool              `json:"is_hazardous"`
	Weight       float64           `json:"weight"`
	Length       float64           `json:"length"`
	Width        float64           `json:"width"`
	Height       float64           `json:"height"`
	PackagedSize float64           `json:"packaged_size"`
	IsNew        bool              `json:"is_new"`
	IsActive     bool              `json:"is_active"`
	Rating       float64           `json:"rating"`
	ShopID       int               `json:"shop_id"`
	CategoryId   int               `json:"category_id"`
	BulkPrice    *ProductBulkPrice `json:"bulkPrice"`
	VariantGroup []*VariantGroup   `json:"variantGroup"`
	Media        []*ProductMedia   `json:"media"`

	gorm.Model `json:"-"`
}
