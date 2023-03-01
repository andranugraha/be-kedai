package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type ShopVoucherRepository interface {
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
}

type shopVoucherRepositoryImpl struct {
	db *gorm.DB
}

type ShopVoucherRConfig struct {
	DB *gorm.DB
}

func NewShopVoucherRepository(cfg *ShopVoucherRConfig) ShopVoucherRepository {
	return &shopVoucherRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *shopVoucherRepositoryImpl) GetShopVoucher(shopId int) ([]*model.ShopVoucher, error) {
	var shopVoucher []*model.ShopVoucher
	publicVoucher := true
	err := r.db.Where("shop_id = ?", shopId).Where("is_hidden != ?", publicVoucher).Where("now() < expired_at").Find(&shopVoucher).Error
	if err != nil {
		return nil, err
	}

	return shopVoucher, nil
}
