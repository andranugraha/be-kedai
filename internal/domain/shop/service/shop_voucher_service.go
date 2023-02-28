package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopVoucherService interface {
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
}

type shopVoucherServiceImpl struct {
	shopVoucherRepository repository.ShopVoucherRepository
}

type ShopVoucherSConfig struct {
	ShopVoucherRepository repository.ShopVoucherRepository
}

func NewShopVoucherService(cfg *ShopVoucherSConfig) ShopVoucherService {
	return &shopVoucherServiceImpl {
		shopVoucherRepository: cfg.ShopVoucherRepository,
	}
}

func (s *shopVoucherServiceImpl) GetShopVoucher(shopId int) ([]*model.ShopVoucher, error) {
	return s.shopVoucherRepository.GetShopVoucher(shopId)
}