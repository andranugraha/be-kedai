package repository

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"time"

	"gorm.io/gorm"
)

type MarketplaceBannerRepository interface {
	GetMarketplaceBanner() ([]*model.MarketplaceBanner, error)
}

type marketplaceBannerRepositoryImpl struct {
	db *gorm.DB
}

type MarketplaceBannerRConfig struct {
	DB *gorm.DB
}

func NewMarketplaceBannerRepository(cfg *MarketplaceBannerRConfig) MarketplaceBannerRepository {
	return &marketplaceBannerRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *marketplaceBannerRepositoryImpl) GetMarketplaceBanner() ([]*model.MarketplaceBanner, error) {
	var banners []*model.MarketplaceBanner

	currentTime := time.Now()

	err := r.db.
		Model(&model.MarketplaceBanner{}).
		Where("? BETWEEN start_date AND end_date", currentTime).
		Find(&banners).Error
	if err != nil {
		return nil, err
	}

	return banners, nil
}
