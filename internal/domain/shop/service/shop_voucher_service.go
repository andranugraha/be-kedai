package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopVoucherService interface {
	GetValidShopVoucherById(id int) (*model.ShopVoucher, error)
	GetShopVoucher(slug string) ([]*model.ShopVoucher, error)
	GetValidShopVoucherByUserIDAndSlug(userID int, slug string) ([]*model.ShopVoucher, error)
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

func (s *shopVoucherServiceImpl) GetValidShopVoucherById(id int) (*model.ShopVoucher, error) {
	return s.shopVoucherRepository.GetValidById(id)
}

func (s *shopVoucherServiceImpl) GetShopVoucher(slug string) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetShopVoucher(shop.ID)
}

func (s *shopVoucherServiceImpl) GetValidShopVoucherByUserIDAndSlug(userID int, slug string) ([]*model.ShopVoucher, error) {
	shop, err := s.shopService.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.shopVoucherRepository.GetValidByUserIDAndShopID(userID, shop.ID)
}
