package model

import (
	"time"

	"gorm.io/gorm"
)

type RefundRequest struct {
	ID           int       `json:"id"`
	RequestDate  time.Time `json:"requestDate"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	RefundAmount float64   `json:"refundAmount"`
	InvoiceId    int       `json:"invoiceId"`

	Invoice *InvoicePerShop `json:"invoice"`

	gorm.Model `json:"-"`
}
