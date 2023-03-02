package repository

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"time"

	"gorm.io/gorm"
)

type MarketplaceVoucherRepository interface {
	GetMarketplaceVoucher() ([]*model.MarketplaceVoucher, error)
}

type marketplaceVoucherRepositoryImpl struct {
	db *gorm.DB
}

type MarketplaceVoucherRConfig struct {
	DB *gorm.DB
}

func NewMarketplaceVoucherRepository(cfg *MarketplaceVoucherRConfig) MarketplaceVoucherRepository {
	return &marketplaceVoucherRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *marketplaceVoucherRepositoryImpl) GetMarketplaceVoucher() ([]*model.MarketplaceVoucher, error) {
	var marketplaceVoucher []*model.MarketplaceVoucher

	publicVoucher := true
	err := r.db.Where("expired_at > ?", time.Now()).Where("is_hidden != ?", publicVoucher).Find(&marketplaceVoucher).Error
	if err != nil {
		return nil, err
	}

	return marketplaceVoucher, nil
}
