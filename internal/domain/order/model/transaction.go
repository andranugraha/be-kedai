package model

import "gorm.io/gorm"

type Transaction struct {
	ID         int     `json:"id"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"totalPrice"`
	Quantity   int     `json:"quantity"`
	Note       *string `json:"note"`

	InvoiceID int `json:"invoiceId"`
	AddressID int `json:"addressId"`
	UserID    int `json:"userId"`
	SkuID     int `json:"skuId"`

	gorm.Model `json:"-"`
}
