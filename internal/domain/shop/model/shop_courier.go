package model

import "gorm.io/gorm"

type ShopCourier struct {
	ID               int
	ShopID           int
	CourierServiceID int
	IsActive         bool

	Shop           *Shop
	CourierService *CourierService

	gorm.Model
}
