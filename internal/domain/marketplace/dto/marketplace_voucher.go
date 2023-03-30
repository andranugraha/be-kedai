package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
)

type GetMarketplaceVoucherRequest struct {
	UserId          int
	CategoryId      int    `form:"categoryId"`
	PaymentMethodId int    `form:"paymentMethodId"`
	Code            string `form:"code"`
}

func (r *GetMarketplaceVoucherRequest) Validate() {

	if r.CategoryId < 0 {
		r.CategoryId = 0
	}

	if r.PaymentMethodId < 0 {
		r.PaymentMethodId = 0
	}

}

type AdminMarketplaceVoucher struct {
	model.MarketplaceVoucher
	Status string `json:"status" gorm:"column:status"`
}

func (AdminMarketplaceVoucher) TableName() string {
	return "marketplace_vouchers"
}

type AdminVoucherFilterRequest struct {
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
	Status string `form:"status"`
	Name   string `form:"name"`
	Code   string `form:"code"`
}

func (p *AdminVoucherFilterRequest) Validate() {
	if p.Limit < 1 {
		p.Limit = constant.DefaultSellerVoucherLimit
	}

	if p.Limit > 50 {
		p.Limit = constant.MaxSellerVoucherLimit
	}

	if p.Page < 1 {
		p.Page = 1
	}
}

func (p *AdminVoucherFilterRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
