package service

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
	productUtils "kedai/backend/be-kedai/internal/utils/product"
)

type MarketplaceVoucherService interface {
	GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetMarketplaceVoucherAdminByCode(voucherCode string) (*dto.AdminMarketplaceVoucher, error)
	GetMarketplaceVoucherAdmin(request *dto.AdminVoucherFilterRequest) (*commonDto.PaginationResponse, error)
	GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetValidForCheckout(id, userID, PaymentMethodID int) (*model.MarketplaceVoucher, error)
	UpdateVoucher(voucherCode string, request *dto.UpdateVoucherRequest) error
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

func (s *marketplaceVoucherServiceImpl) GetMarketplaceVoucherAdminByCode(voucherCode string) (*dto.AdminMarketplaceVoucher, error) {
	return s.marketplaceVoucherRepository.GetMarketplaceVoucherAdminByCode(voucherCode)
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

func (s *marketplaceVoucherServiceImpl) UpdateVoucher(voucherCode string, request *dto.UpdateVoucherRequest) error {
	if isVoucherNameValid := productUtils.ValidateProductName(request.Name); !isVoucherNameValid {
		return commonErr.ErrInvalidVoucherNamePattern
	}

	voucher, err := s.marketplaceVoucherRepository.GetMarketplaceVoucherAdminByCode(voucherCode)
	if err != nil {
		return err
	}

	if voucher.Status == constant.VoucherPromotionStatusExpired {
		return commonErr.ErrVoucherStatusConflict
	}

	isZero := 0

	if request.Name == "" {
		request.Name = voucher.Name
	}
	if request.IsHidden == nil {
		request.IsHidden = &voucher.IsHidden
	}
	if request.Description == "" {
		request.Description = voucher.Description
	}
	if err := request.ValidateDateRange(voucher.ExpiredAt); err != nil {
		return err
	}
	if request.CategoryId == &isZero {
		request.CategoryId = voucher.CategoryID
	}
	if request.PaymentMethodId == &isZero {
		request.PaymentMethodId = voucher.PaymentMethodID
	}

	payload := &model.MarketplaceVoucher{
		ID:              voucher.ID,
		Name:            request.Name,
		Code:            voucher.Code,
		IsHidden:        *request.IsHidden,
		Description:     request.Description,
		ExpiredAt:       request.ExpiredAt,
		CategoryID:      request.CategoryId,
		PaymentMethodID: request.PaymentMethodId,
	}

	err = s.marketplaceVoucherRepository.Update(payload)
	if err != nil {
		return err
	}

	return nil
}
