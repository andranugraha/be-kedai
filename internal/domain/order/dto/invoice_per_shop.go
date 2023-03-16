package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"time"
)

type InvoicePerShopDetail struct {
	model.InvoicePerShop
	MarketplaceVoucherAmount float64            `json:"marketplaceVoucherAmount"`
	MarketplaceVoucherType   string             `json:"marketplaceVoucherType"`
	PaymentDate              time.Time          `json:"paymentDate"`
	TransactionItems         []*TransactionItem `json:"transactionItems" gorm:"foreignKey:InvoiceID"`
}

func (InvoicePerShopDetail) TableName() string {
	return "invoice_per_shops"
}

type TransactionItem struct {
	model.Transaction
	ProductName string            `json:"productName"`
	ImageUrl    string            `json:"imageUrl"`
	Sku         *productModel.Sku `json:"sku"`
}

func (TransactionItem) TableName() string {
	return "transactions"
}

type InvoicePerShopFilterRequest struct {
	S              string `form:"s"`
	Username       string `form:"user"`
	ProductName    string `form:"product"`
	TrackingNumber string `form:"track"`
	OrderId        string `form:"orderId"`
	Limit          int    `form:"limit"`
	Page           int    `form:"page"`
	StartDate      string `form:"startDate" binding:"required_with=EndDate,omitempty,datetime=2006-01-02"`
	EndDate        string `form:"endDate" binding:"required_with=StartDate,omitempty,datetime=2006-01-02"`
	Status         string `form:"status"`
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
		d.Status != constant.TransactionStatusCanceled &&
		d.Status != constant.Released &&
		d.Status != constant.ToRelease {
		d.Status = ""
	}
}

func (d *InvoicePerShopFilterRequest) Offset() int {
	return (d.Page - 1) * d.Limit
}

type WithdrawInvoiceRequest struct {
	OrderID []int `json:"orderId" binding:"required,min=1"`
}

func (d *WithdrawInvoiceRequest) Validate() {
	if len(d.OrderID) == 0 {
		d.OrderID = []int{}
	}
}
