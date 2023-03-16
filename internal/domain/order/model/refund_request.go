package model

import (
	"time"

	"gorm.io/gorm"
)

type RefundRequest struct {
	ID          int       `json:"id"`
	RequestDate time.Time `json:"requestDate"`
	Status      string    `json:"status"`
	InvoiceId   int       `json:"invoiceId"`

	Invoice *InvoicePerShop `json:"invoice"`

	gorm.Model `json:"-"`
}
