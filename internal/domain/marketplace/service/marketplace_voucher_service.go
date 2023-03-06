package service

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
)

type MarketplaceVoucherService interface {
	GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
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

func (s *marketplaceVoucherServiceImpl) GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error) {
	return s.marketplaceVoucherRepository.GetMarketplaceVoucher(req)
}

func (s *marketplaceVoucherServiceImpl) GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error) {
	return s.marketplaceVoucherRepository.GetValidByUserID(req)
}
