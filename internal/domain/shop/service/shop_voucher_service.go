package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopVoucherService interface {
	GetValidShopVoucherByIdAndUserId(id, userId int) (*model.ShopVoucher, error)
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

func (s *shopVoucherServiceImpl) GetSellerVoucher(slug string) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetShopVoucher(shop.ID)
}


func (s *shopVoucherServiceImpl) GetValidShopVoucherByUserIDAndSlug(req dto.GetValidShopVoucherRequest) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(req.Slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetValidByUserIDAndShopID(req, shop.ID)
}
