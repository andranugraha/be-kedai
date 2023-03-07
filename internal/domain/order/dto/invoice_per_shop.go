package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"time"
)

type InvoicePerShopDetail struct {
	model.InvoicePerShop
	VoucherAmount    float64            `json:"voucherAmount,omitempty"`
	VoucherType      string             `json:"voucherType,omitempty"`
	PaymentDate      time.Time          `json:"paymentDate"`
	TransactionItems []*TransactionItem `json:"transactionItems" gorm:"foreignKey:InvoiceID"`
}

func (InvoicePerShopDetail) TableName() string {
	return "invoice_per_shops"
}

type TransactionItem struct {
	model.Transactions
	ProductName string            `json:"productName"`
	ImageUrl    string            `json:"imageUrl"`
	Sku         *productModel.Sku `json:"sku"`
}

func (TransactionItem) TableName() string {
	return "transactions"
}

type InvoicePerShopFilterRequest struct {
	S         string `form:"s"`
	Limit     int    `form:"limit"`
	Page      int    `form:"page"`
	StartDate string `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
	Status    string `form:"status"`
}

func (d *InvoicePerShopFilterRequest) Validate() {
	if d.Limit < 1 {
		d.Limit = 10
	}

	if d.Page < 1 {
		d.Page = 1
	}

	if d.Status != constant.TransactionStatusComplained &&
		d.Status != constant.TransactionStatusComplaintConfirmed &&
		d.Status != constant.TransactionStatusComplaintRejected &&
		d.Status != constant.TransactionStatusCompleted &&
		d.Status != constant.TransactionStatusCreated &&
		d.Status != constant.TransactionStatusReceived &&
		d.Status != constant.TransactionStatusSent &&
		d.Status != constant.TransactionStatusCancelled {
		d.Status = ""
	}
}

func (d *InvoicePerShopFilterRequest) Offset() int {
	return (d.Page - 1) * d.Limit
}
