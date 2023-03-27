package service

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
)

type MarketplaceBannerService interface {
	GetMarketplaceBanner() ([]*model.MarketplaceBanner, error)
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
