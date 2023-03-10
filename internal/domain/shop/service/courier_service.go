package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type CourierService interface {
	GetAllCouriers() ([]*model.Courier, error)
	GetShipmentList(userId int) ([]*dto.ShipmentCourierResponse, error)
	GetCouriersByShopID(shopID int) ([]*model.Courier, error)
	GetCourierByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error)
	GetCouriersByProductID(productID int) ([]*model.Courier, error)
	GetMatchingCouriersByShopIDAndProductIDs(*dto.MatchingProductCourierRequest) ([]*model.Courier, error)
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

func (s *courierServiceImpl) GetShipmentList(userId int) ([]*dto.ShipmentCourierResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	return s.courierRepository.GetShipmentList(shop.ID)
}

func (s *courierServiceImpl) GetCouriersByShopID(shopID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByShopID(shopID)
}

func (s *courierServiceImpl) GetCourierByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error) {
	return s.courierRepository.GetByServiceIDAndShopID(courierID, shopID)
}

func (s *courierServiceImpl) GetCouriersByProductID(productID int) ([]*model.Courier, error) {
	return s.courierRepository.GetByProductID(productID)
}

func (s *courierServiceImpl) GetMatchingCouriersByShopIDAndProductIDs(req *dto.MatchingProductCourierRequest) ([]*model.Courier, error) {
	shop, err := s.shopService.FindShopBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	req.ShopID = shop.ID

	return s.courierRepository.GetMatchingCouriersByShopIDAndProductIDs(req)
}
