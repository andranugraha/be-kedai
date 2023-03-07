package model

import (
	"fmt"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID            int       `json:"id"`
	Code          string    `json:"code"`
	Total         float64   `json:"total"`
	Subtotal      float64   `json:"subtotal"`
	VoucherAmount *float64  `json:"voucherAmount,omitempty"`
	VoucherType   *string   `json:"voucherType,omitempty"`
	PaymentDate   time.Time `json:"paymentDate" gorm:"default:CURRENT_TIMESTAMP"`

	UserID          int  `json:"userId"`
	VoucherID       *int `json:"voucherId,omitempty"`
	PaymentMethodID int  `json:"paymentMethodId"`
	UserAddressID   int  `json:"userAddressId"`

	Voucher         *userModel.UserVoucher `json:"voucher,omitempty" gorm:"foreignKey:VoucherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InvoicePerShops []*InvoicePerShop      `json:"invoicePerShops" gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	gorm.Model `json:"-"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	var currentTotal int64
	tx.Model(&Invoice{}).Count(&currentTotal)

	now := time.Now()
	i.Code = fmt.Sprintf("INV/%d%d%d/%d", now.Year(), now.Month(), now.Day(), currentTotal+1)

	return
}
