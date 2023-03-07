package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type CourierService interface {
	GetAllCouriers() ([]*model.Courier, error)
	GetCouriersByShopID(shopID int) ([]*model.Courier, error)
	GetCouriersByProductID(productID int) ([]*model.Courier, error)
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

func (s *courierServiceImpl) GetAllCouriers() ([]*model.Courier, error) {
	return s.courierRepository.GetAll()
}

func (s *courierServiceImpl) GetCouriersByShopID(shopID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByShopID(shopID)
}

func (s *courierServiceImpl) GetCouriersByProductID(productID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByProductID(productID)
}
