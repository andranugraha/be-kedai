package model

import (
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type Product struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	View        int     `json:"view"`
	IsHazardous bool    `json:"isHazardous"`
	Weight      float64 `json:"weight"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	IsNew       bool    `json:"isNew"`
	IsActive    bool    `json:"isActive"`
	Rating      float64 `json:"rating"`
	Sold        int     `json:"sold"`

	ShopID         int                         `json:"shopId"`
	Shop           *shopModel.Shop             `json:"shop,omitempty"`
	CategoryID     int                         `json:"categoryId"`
	Bulk           *ProductBulkPrice           `json:"bulk,omitempty"`
	VariantGroup   []*VariantGroup             `json:"variantGroups,omitempty"`
	Media          []*ProductMedia             `json:"media,omitempty"`
	SKU            *Sku                        `json:"sku,omitempty"`
	SKUs           []*Sku                      `json:"skus,omitempty"`
	CourierService []*shopModel.CourierService `json:"courierServices,omitempty" gorm:"many2many:product_couriers"`

	gorm.Model `json:"-"`
}
