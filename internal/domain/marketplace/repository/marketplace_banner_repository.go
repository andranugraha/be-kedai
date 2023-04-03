package repository

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/utils/date"
	"time"

	"gorm.io/gorm"
)

type MarketplaceBannerRepository interface {
	GetMarketplaceBanner() ([]*model.MarketplaceBanner, error)
	AddMarketplaceBanner(body *dto.MarketplaceBannerRequest) (*model.MarketplaceBanner, error)
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

func (r *marketplaceBannerRepositoryImpl) Last(banner *model.MarketplaceBanner) (*model.MarketplaceBanner, error) {
	result := r.db.Last(&banner)
	return banner, result.Error
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

func (r *marketplaceBannerRepositoryImpl) AddMarketplaceBanner(body *dto.MarketplaceBannerRequest) (*model.MarketplaceBanner, error) {
	newBanner := &model.MarketplaceBanner{
		MediaUrl:  body.MediaUrl,
		StartDate: date.ParseRFC3999NanoTime(body.StartDate, time.Now()),
		EndDate:   date.ParseRFC3999NanoTime(body.EndDate, time.Now().AddDate(0, 0, 14)),
	}

	result := r.db.Create(&newBanner)

	if result.Error != nil {
		return nil, result.Error
	}

	banner, _ := r.Last(newBanner)

	return banner, nil
}
