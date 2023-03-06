package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type CourierService interface {
	GetCouriersByShopID(shopID int) ([]*model.Courier, error)
	GetCourierByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error)
}

type courierServiceImpl struct {
	courierRepository repository.CourierRepository
	shopService       ShopService
}

type CourierSConfig struct {
	CourierRepository repository.CourierRepository
	ShopService       ShopService
}

func NewCourierService(cfg *CourierSConfig) CourierService {
	return &courierServiceImpl{
		courierRepository: cfg.CourierRepository,
		shopService:       cfg.ShopService,
	}
}

func (s *courierServiceImpl) GetCouriersByShopID(shopID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByShopID(shopID)
}

func (s *courierServiceImpl) GetCourierByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error) {
	courier, err := s.courierRepository.GetByServiceIDAndShopID(courierID, shopID)
	if err != nil {
		return nil, err
	}

	return courier, nil
}
