package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	"log"
)

type ShopService interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
}

type shopServiceImpl struct {
	shopRepository repository.ShopRepository
	shopVoucherService ShopVoucherService
}

type ShopSConfig struct {
	ShopRepository repository.ShopRepository
	ShopVoucherService ShopVoucherService
}

func NewShopService(cfg *ShopSConfig) ShopService {
	return &shopServiceImpl{
		shopRepository: cfg.ShopRepository,
		shopVoucherService: cfg.ShopVoucherService,
	}
}

func (s *shopServiceImpl) FindShopById(id int) (*model.Shop, error) {
	return s.shopRepository.FindShopById(id)
}

func (s *shopServiceImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	return s.shopRepository.FindShopByUserId(userId)
}

func (s *shopServiceImpl) FindShopBySlug(slug string) (*model.Shop, error) {
	shop, err := s.shopRepository.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}
	
	voucher, err := s.shopVoucherService.GetShopVoucher(shop.ID)
	if err != nil {
		return nil, err
	}

	log.Print(voucher)

	shop.ShopVoucher = voucher

	return shop, nil
}
