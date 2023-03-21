package dto

import "kedai/backend/be-kedai/internal/domain/shop/model"

type GetValidShopVoucherRequest struct {
	Slug   string
	Code   string
	UserID int
}

type SellerVoucher struct {
	model.ShopVoucher
	Status string `json:"status" gorm:"column:status"`
}

func (SellerVoucher) TableName() string {
	return "shop_vouchers"
}

type SellerVoucherFilterRequest struct {
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
	Status string `form:"status"`
	Name   string `form:"name"`
	Code   string `form:"code"`
}

func (p *SellerVoucherFilterRequest) Validate() {
	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}
}

func (p *SellerVoucherFilterRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
