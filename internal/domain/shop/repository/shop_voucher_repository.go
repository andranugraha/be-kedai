package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type ShopVoucherRepository interface {
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
	GetValidById(id int) (*model.ShopVoucher, error)
}

type shopVoucherImpl struct {
	db *gorm.DB
}

type ShopVoucherRConfig struct {
	DB *gorm.DB
}

func NewShopVoucherRepository(cfg *ShopVoucherRConfig) ShopVoucherRepository {
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

func (r *shopVoucherImpl) GetValidById(id int) (*model.ShopVoucher, error) {
	var shopVoucher model.ShopVoucher

	err := r.db.Where("expired_at > now()").First(&shopVoucher, id).Error
	if err != nil {
		return nil, err
	}

	return &shopVoucher, nil
}
