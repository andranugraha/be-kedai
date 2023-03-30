package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
)

type MarketplaceVoucherService interface {
	GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetMarketplaceVoucherAdmin(request *dto.AdminVoucherFilterRequest) (*commonDto.PaginationResponse, error)
	GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetValidForCheckout(id, userID, PaymentMethodID int) (*model.MarketplaceVoucher, error)
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

func (s *marketplaceVoucherServiceImpl) GetMarketplaceVoucherAdmin(req *dto.AdminVoucherFilterRequest) (*commonDto.PaginationResponse, error) {
	vouchers, totalRows, totalPages, err := s.marketplaceVoucherRepository.GetMarketplaceVoucherAdmin(req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Page:       req.Page,
		Limit:      req.Limit,
		Data:       vouchers,
	}, nil
}

func (s *marketplaceVoucherServiceImpl) GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error) {
	return s.marketplaceVoucherRepository.GetValidByUserID(req)
}

func (s *marketplaceVoucherServiceImpl) GetValidForCheckout(id, userID, PaymentMethodID int) (*model.MarketplaceVoucher, error) {
	return s.marketplaceVoucherRepository.GetValid(id, userID, PaymentMethodID)
}
