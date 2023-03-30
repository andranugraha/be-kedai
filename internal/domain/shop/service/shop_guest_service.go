package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopGuestService interface {
	CreateShopGuest(shopId int) (*model.ShopGuest, error)
}

type shopGuestServiceImpl struct {
	shopGuestRepository repository.ShopGuestRepository
}

type ShopGuestSConfig struct {
	ShopGuestRepository repository.ShopGuestRepository
}

func NewShopGuestService(cfg *ShopGuestSConfig) ShopGuestService {
	return &shopGuestServiceImpl{
		shopGuestRepository: cfg.ShopGuestRepository,
	}
}

func (s *shopGuestServiceImpl) CreateShopGuest(shopId int) (*model.ShopGuest, error) {
	return s.shopGuestRepository.CreateShopGuest(&model.ShopGuest{ShopId: shopId})
}
