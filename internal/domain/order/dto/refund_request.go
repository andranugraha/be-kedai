package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
)

type RefundRequest struct {
	RefundStatus string `json:"refundStatus" binding:"required"`
}

func (c *RefundRequest) Validate() error {
	if (c.RefundStatus != constant.RequestStatusSellerApproved) && (c.RefundStatus != constant.RefundStatusRejected) {
		return commonErr.ErrInvalidRefundStatus
	}
	return nil
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
	Quantity             int     `gorm:"column:quantity" json:"quantity"`
}

type GetRefund struct {
	RequestRefundId     int     `gorm:"column:id" json:"requestRefundId"`
	RequestRefundStatus string  `gorm:"column:status" json:"requestRefundStatus"`
	RequestRefundType   string  `gorm:"column:type" json:"requestRefundType"`
	InvoicePerShopId    int     `gorm:"column:id" json:"invoicePerShopId"`
	RefundAmount        float64 `gorm:"column:refund_amount" json:"refundAmount"`
	InvoiceCode         string  `gorm:"column:code" json:"invoiceCode"`
	InvoiceTotal        float64 `gorm:"column:total" json:"invoiceTotal"`
	ShippingCost        float64 `gorm:"column:shipping_cost" json:"shippingCost"`
	ProductName         string  `gorm:"column:name" json:"productName"`
	ProductMedia        string  `gorm:"column:url" json:"productMedia"`
	Username            string  `gorm:"column:username" json:"username"`
}

type GetRefundReq struct {
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Status string `json:"status"`
}

func (g *GetRefundReq) Validate() {

	if g.Limit < 1 {
		g.Limit = constant.MinLimitGetRefundRequest
	}

	if g.Limit > constant.MaxLimitGetRefundRequest {
		g.Limit = constant.MaxLimitGetRefundRequest
	}

	if g.Page < 1 {
		g.Page = 1
	}
}
