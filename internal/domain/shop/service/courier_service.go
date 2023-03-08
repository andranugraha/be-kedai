package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type CourierService interface {
	GetShipmentList(shopId int) ([]*dto.ShipmentCourierResponse, error)
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

func (s *courierServiceImpl) GetShipmentList(shopId int) ([]*dto.ShipmentCourierResponse, error) {
	return s.courierRepository.GetShipmentList(shopId)
}

func (s *courierServiceImpl) GetCouriersByShopID(shopID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByShopID(shopID)
}

func (s *courierServiceImpl) GetCouriersByProductID(productID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByProductID(productID)
}
