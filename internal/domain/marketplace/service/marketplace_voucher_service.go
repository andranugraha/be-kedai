package service

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
)

type MarketplaceVoucherService interface {
	GetMarketplaceVoucher() ([]*model.MarketplaceVoucher, error)
}

type marketplaceVoucherServiceImpl struct {
	marketplaceVoucherRepository repository.MarketplaceVoucherRepository
}

type MarketplaceVoucherSConfig struct {
	MarketplaceVoucherRepository repository.MarketplaceVoucherRepository
}

func NewMarketplaceVoucherService(cfg *MarketplaceVoucherSConfig) MarketplaceVoucherService {
	return &marketplaceVoucherServiceImpl{
		marketplaceVoucherRepository: cfg.MarketplaceVoucherRepository,
	}
}

func (s *marketplaceVoucherServiceImpl) GetMarketplaceVoucher() ([]*model.MarketplaceVoucher, error) {
	return s.marketplaceVoucherRepository.GetMarketplaceVoucher()
}
