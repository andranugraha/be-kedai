package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	productUtils "kedai/backend/be-kedai/internal/utils/product"
)

type ShopVoucherService interface {
	GetValidShopVoucherByIdAndUserId(id, userId int) (*model.ShopVoucher, error)
	GetSellerVoucher(userID int, req *dto.SellerVoucherFilterRequest) (*commonDto.PaginationResponse, error)
	GetShopVoucher(slug string) ([]*model.ShopVoucher, error)
	GetValidShopVoucherByUserIDAndSlug(dto.GetValidShopVoucherRequest) ([]*model.ShopVoucher, error)
}

type shopVoucherServiceImpl struct {
	shopVoucherRepository repository.ShopVoucherRepository
	shopService           ShopService
}

type ShopVoucherSConfig struct {
	ShopVoucherRepository repository.ShopVoucherRepository
	ShopService           ShopService
}

func NewShopVoucherService(cfg *ShopVoucherSConfig) ShopVoucherService {
	return &shopVoucherServiceImpl{
		shopVoucherRepository: cfg.ShopVoucherRepository,
		shopService:           cfg.ShopService,
	}
}

func (s *shopVoucherServiceImpl) GetValidShopVoucherByIdAndUserId(id, userId int) (*model.ShopVoucher, error) {
	return s.shopVoucherRepository.GetValidByIdAndUserId(id, userId)
}

func (s *shopVoucherServiceImpl) GetShopVoucher(slug string) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetShopVoucher(shop.ID)
}

func (s *shopVoucherServiceImpl) GetSellerVoucher(userID int, req *dto.SellerVoucherFilterRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	vouchers, totalRows, totalPages, err := s.shopVoucherRepository.GetSellerVoucher(shop.ID, req)
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

func (s *shopVoucherServiceImpl) CreateVoucher(userID int, request *dto.CreateVoucherRequest) (*model.ShopVoucher, error) {
	if isProductNameValid := productUtils.ValidateProductName(request.Name); !isProductNameValid {
		return nil, commonErr.ErrInvalidProductNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	voucher, err := s.shopVoucherRepository.Create(shop.ID, request)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *shopVoucherServiceImpl) GetValidShopVoucherByUserIDAndSlug(req dto.GetValidShopVoucherRequest) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(req.Slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetValidByUserIDAndShopID(req, shop.ID)
}
