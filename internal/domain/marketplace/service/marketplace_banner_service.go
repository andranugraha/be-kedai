package service

import (
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
	"kedai/backend/be-kedai/internal/utils/date"
	"time"
)

type MarketplaceBannerService interface {
	GetMarketplaceBanner() ([]*model.MarketplaceBanner, error)
	AddMarketplaceBanner(body *dto.MarketplaceBannerRequest) (*model.MarketplaceBanner, error)
}

type marketplaceBannerServiceImpl struct {
	marketplaceBannerRepository repository.MarketplaceBannerRepository
}

type MarketplaceBannerSConfig struct {
	MarketplaceBannerRepository repository.MarketplaceBannerRepository
}

func NewMarketplaceBannerService(cfg *MarketplaceBannerSConfig) MarketplaceBannerService {
	return &marketplaceBannerServiceImpl{
		marketplaceBannerRepository: cfg.MarketplaceBannerRepository,
	}
}

func (s *marketplaceBannerServiceImpl) GetMarketplaceBanner() ([]*model.MarketplaceBanner, error) {
	return s.marketplaceBannerRepository.GetMarketplaceBanner()
}

func (s *marketplaceBannerServiceImpl) AddMarketplaceBanner(body *dto.MarketplaceBannerRequest) (*model.MarketplaceBanner, error) {
	if !date.IsValidRFC3999NanoDate(body.StartDate) || !date.IsValidRFC3999NanoDate(body.EndDate) {
		return nil, spErr.ErrInvalidRFC3999Nano
	}
	if !date.ParseRFC3999NanoTime(body.StartDate, time.Now()).Before(date.ParseRFC3999NanoTime(body.EndDate, time.Now().AddDate(0, 0, 14))) {
		return nil, spErr.ErrBackDate
	}
	return s.marketplaceBannerRepository.AddMarketplaceBanner(body)
}
