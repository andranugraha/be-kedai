package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopService interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
}

type shopServiceImpl struct {
	shopRepository repository.ShopRepository
}

type ShopSConfig struct {
	ShopRepository repository.ShopRepository
}

func NewShopService(cfg *ShopSConfig) ShopService {
	return &shopServiceImpl{
		shopRepository: cfg.ShopRepository,
	}
}

func (s *shopServiceImpl) FindShopById(id int) (*model.Shop, error) {
	return s.shopRepository.FindShopById(id)
}

func (s *shopServiceImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	return s.shopRepository.FindShopByUserId(userId)
}
