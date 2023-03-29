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

type RefundInfo struct {
	RequestRefundId      int     `gorm:"column:id" json:"requestRefundId"`
	RequestRefundStatus  string  `gorm:"column:status" json:"requestRefundStatus"`
	RequestRefundType    string  `gorm:"column:type" json:"requestRefundType"`
	InvoicePerShopId     int     `gorm:"column:invoice_id" json:"invoicePerShopId"`
	RefundAmount         float64 `gorm:"column:refund_amount" json:"refundAmount"`
	ShippingCost         float64 `gorm:"column:shipping_cost" json:"shippingCost"`
	ShopVoucherId        int     `gorm:"column:voucher_id" json:"shopVoucherId"`
	ShopId               int     `gorm:"column:shop_id" json:"shopId"`
	UserId               int     `gorm:"column:user_id" json:"userId"`
	WalletId             int     `gorm:"column:id" json:"walletId"`
	InvoiceId            int     `gorm:"column:id" json:"invoiceId"`
	MarketplaceVoucherId int     `gorm:"column:voucher_id" json:"marketplaceVoucherId"`
	SkuId                int     `gorm:"column:sku_id" json:"skuId"`
}
