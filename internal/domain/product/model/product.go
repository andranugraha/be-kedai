package model

import (
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type Product struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	View        int     `json:"view,omitempty"`
	IsHazardous bool    `json:"isHazardous,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
	Length      float64 `json:"length,omitempty"`
	Width       float64 `json:"width,omitempty"`
	Height      float64 `json:"height,omitempty"`
	IsNew       bool    `json:"isNew,omitempty"`
	IsActive    bool    `json:"isActive,omitempty"`
	Rating      float64 `json:"rating,omitempty"`
	Sold        int     `json:"sold,omitempty"`

	ShopID         int                         `json:"shopId,omitempty"`
	Shop           *shopModel.Shop             `json:"shop,omitempty"`
	CategoryID     int                         `json:"categoryId,omitempty"`
	Bulk           *ProductBulkPrice           `json:"bulk,omitempty"`
	VariantGroup   []*VariantGroup             `json:"variantGroups,omitempty"`
	Media          []*ProductMedia             `json:"media,omitempty"`
	SKU            *Sku                        `json:"sku,omitempty"`
	SKUs           []*Sku                      `json:"skus,omitempty"`
	CourierService []*shopModel.CourierService `json:"courierServices,omitempty" gorm:"many2many:product_couriers"`

	gorm.Model `json:"-"`
}
