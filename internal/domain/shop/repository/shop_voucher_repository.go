package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type ShopVoucherRepository interface{
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
}

type shopVoucherImpl struct {
	db *gorm.DB
}

type ShopVConfig struct {
	DB *gorm.DB
}

func NewShopVoucherRepository(cfg *ShopVConfig) ShopVoucherRepository {
	return &shopVoucherImpl{
		db: cfg.DB,
	}
}

func (r *shopVoucherImpl) GetShopVoucher(shopId int) ([]*model.ShopVoucher, error) {
	var shopVoucher []*model.ShopVoucher

	err := r.db.Where("shop_id = ?", shopId).Find(&shopVoucher).Error
	if err != nil {
		return nil, err
	}

	return shopVoucher, nil
}
