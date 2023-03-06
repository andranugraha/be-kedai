package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"time"

	"gorm.io/gorm"
)

type ShopVoucherRepository interface {
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
	GetValidById(id int) (*model.ShopVoucher, error)
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
	now := time.Now()
	err := r.db.Where("shop_id = ?", shopId).Where("is_hidden != ?", publicVoucher).Where("? < expired_at", now).Find(&shopVoucher).Error
	if err != nil {
		return nil, err
	}

	return shopVoucher, nil
}

func (r *shopVoucherRepositoryImpl) GetValidById(id int) (*model.ShopVoucher, error) {
	var shopVoucher model.ShopVoucher

	err := r.db.Where("expired_at > now()").First(&shopVoucher, id).Error
	if err != nil {
		return nil, err
	}

	return &shopVoucher, nil
}
