package dto

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"time"
)

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

type CreateVoucherRequest struct {
	Name         string    `json:"name" binding:"required,min=1,max=100"`
	Code         string    `json:"code" binding:"required,alphanum,min=1,max=9"`
	Amount       float64   `json:"amount" binding:"required,min=0,max=500000000"`
	Type         string    `json:"type" binding:"required"`
	IsHidden     *bool     `json:"isHidden" binding:"required"`
	Description  string    `json:"description" binding:"required,min=5,max=1000"`
	MinimumSpend float64   `json:"minimumSpend" binding:"required,min=0,max=500000000"`
	TotalQuota   int       `json:"totalQuota" binding:"required,min=1,max=200000"`
	StartFrom    time.Time `json:"startFrom" binding:"required"`
	ExpiredAt    time.Time `json:"expiredAt" binding:"required"`
}
