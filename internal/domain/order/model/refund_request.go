package model

import (
	"time"

	"gorm.io/gorm"
)

type RefundRequest struct {
	ID            int       `json:"id"`
	RequestDate   time.Time `json:"requestDate"`
	Status        string    `json:"status"`
	TransactionId int       `json:"transactionId"`

	Transaction *Transaction `json:"transaction"`

	gorm.Model `json:"-"`
}
