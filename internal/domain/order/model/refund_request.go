package model

import (
	"time"

	"gorm.io/gorm"
)

type RefundRequest struct {
	ID           int       `json:"id"`
	RequestDate  time.Time `json:"requestDate" gorm:"default:CURRENT_TIMESTAMP"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	RefundAmount float64   `json:"refundAmount"`
	InvoiceID    int       `json:"invoiceId"`

	Invoice *InvoicePerShop `json:"invoice" gorm:"foreignKey:InvoiceID"`

	gorm.Model `json:"-"`
}

func (rr *RefundRequest) BeforeCreate(tx *gorm.DB) (err error) {

	rr.RequestDate = time.Now()
	return
}
