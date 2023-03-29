package dto

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"time"
)

type SellerPromotion struct {
	model.ShopPromotion
	Status  string                                `json:"status" gorm:"column:status"`
	Product []*dto.SellerProductPromotionResponse `json:"products"`
}

func (SellerPromotion) TableName() string {
	return "shop_promotions"
}

type SellerPromotionFilterRequest struct {
	Limit  int    `form:"limit,max=100"`
	Page   int    `form:"page"`
	Status string `form:"status"`
	Name   string `form:"name"`
}

func (p *SellerPromotionFilterRequest) Validate() {
	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}
}

func (p *SellerPromotionFilterRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

type UpdateShopPromotionRequest struct {
	Name              string                               `json:"name" binding:"min=1,max=100"`
	StartPeriod       time.Time                            `json:"startPeriod"`
	EndPeriod         time.Time                            `json:"endPeriod"`
	ProductPromotions []*dto.UpdateProductPromotionRequest `json:"productPromotions"`
}

func (p *UpdateShopPromotionRequest) ValidateDateRange() error {
	now := time.Now().UTC()

	if p.StartPeriod.After(p.EndPeriod) || (p.StartPeriod.Before(now) && p.EndPeriod.Before(now)) || p.EndPeriod.Before(p.StartPeriod) {
		return errs.ErrInvalidVoucherDateRange
	}

	return nil
}
