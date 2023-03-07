package model

import (
	"time"

	"gorm.io/gorm"
)

type InvoiceStatus struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	StatusDate time.Time `json:"statusDate" gorm:"default:CURRENT_TIMESTAMP"`

	InvoicePerShopID int `json:"invoicePerShopId"`

	gorm.Model `json:"-"`
}
