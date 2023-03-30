package service

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
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
	GetVoucherByCodeAndShopId(voucherCode string, userID int) (*dto.SellerVoucher, error)
	CreateVoucher(userID int, request *dto.CreateVoucherRequest) (*model.ShopVoucher, error)
	UpdateVoucher(userID int, voucherCode string, request *dto.UpdateVoucherRequest) (*model.ShopVoucher, error)
	DeleteVoucher(userID int, voucherCode string) error
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

func (s *shopVoucherServiceImpl) GetVoucherByCodeAndShopId(voucherCode string, userID int) (*dto.SellerVoucher, error) {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	voucher, err := s.shopVoucherRepository.GetVoucherByCodeAndShopId(voucherCode, shop.ID)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *shopVoucherServiceImpl) CreateVoucher(userID int, request *dto.CreateVoucherRequest) (*model.ShopVoucher, error) {
	if isVoucherNameValid := productUtils.ValidateProductName(request.Name); !isVoucherNameValid {
		return nil, commonErr.ErrInvalidVoucherNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	existingVoucher, err := s.shopVoucherRepository.GetVoucherByCodeAndShopId(request.Code, shop.ID)
	if err != nil && !errors.Is(err, commonErr.ErrVoucherNotFound) {
		return nil, err
	}

	if existingVoucher != nil && (existingVoucher.Status == constant.VoucherPromotionStatusOngoing || existingVoucher.Status == constant.VoucherPromotionStatusUpcoming) {
		return nil, commonErr.ErrDuplicateVoucherCode
	}

	if err := s.shopVoucherRepository.ValidateVoucherDateRange(request.StartFrom, request.ExpiredAt); err != nil {
		return nil, err
	}

	voucher, err := s.shopVoucherRepository.Create(shop.ID, request)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *shopVoucherServiceImpl) UpdateVoucher(userID int, voucherCode string, request *dto.UpdateVoucherRequest) (*model.ShopVoucher, error) {
	if isVoucherNameValid := productUtils.ValidateProductName(request.Name); !isVoucherNameValid {
		return nil, commonErr.ErrInvalidVoucherNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	voucher, err := s.shopVoucherRepository.GetVoucherByCodeAndShopId(voucherCode, shop.ID)
	if err != nil {
		return nil, err
	}

	if voucher.Status == constant.VoucherPromotionStatusUpcoming {
		if err := s.shopVoucherRepository.ValidateVoucherDateRange(request.StartFrom, request.ExpiredAt); err != nil {
			return nil, err
		}
	} else if voucher.Status == constant.VoucherPromotionStatusOngoing {
		if request.Amount != 0 || request.Type != "" || request.Description != "" || request.MinimumSpend != 0 || !request.StartFrom.IsZero() {
			return nil, commonErr.ErrVoucherFieldsCantBeEdited
		}
	} else {
		return nil, commonErr.ErrVoucherStatusConflict
	}

	if request.Name == "" {
		request.Name = voucher.Name
	}
	if request.Amount == 0 {
		request.Amount = voucher.Amount
	}
	if request.Type == "" {
		request.Type = voucher.Type
	}
	if request.IsHidden == nil {
		request.IsHidden = &voucher.IsHidden
	}
	if request.Description == "" {
		request.Description = voucher.Description
	}
	if request.MinimumSpend == 0 {
		request.MinimumSpend = voucher.MinimumSpend
	}
	if request.TotalQuota == 0 {
		request.TotalQuota = voucher.TotalQuota
	}
	if request.StartFrom.IsZero() {
		request.StartFrom = voucher.StartFrom
	}
	if request.ExpiredAt.IsZero() {
		request.ExpiredAt = voucher.ExpiredAt
	}

	if err := s.shopVoucherRepository.ValidateVoucherDateRange(request.StartFrom, request.ExpiredAt); err != nil {
		return nil, err
	}

	payload := &model.ShopVoucher{
		ID:           voucher.ID,
		Name:         request.Name,
		Code:         voucher.Code,
		Amount:       request.Amount,
		Type:         request.Type,
		IsHidden:     *request.IsHidden,
		Description:  request.Description,
		MinimumSpend: request.MinimumSpend,
		UsedQuota:    voucher.UsedQuota,
		TotalQuota:   request.TotalQuota,
		StartFrom:    request.StartFrom,
		ExpiredAt:    request.ExpiredAt,
		ShopId:       shop.ID,
	}

	res, err := s.shopVoucherRepository.Update(payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *shopVoucherServiceImpl) DeleteVoucher(userID int, voucherCode string) error {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return err
	}

	err = s.shopVoucherRepository.Delete(shop.ID, voucherCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *shopVoucherServiceImpl) GetValidShopVoucherByUserIDAndSlug(req dto.GetValidShopVoucherRequest) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(req.Slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetValidByUserIDAndShopID(req, shop.ID)
}
