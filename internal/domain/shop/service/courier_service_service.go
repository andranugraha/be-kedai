package service

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type CourierServiceService interface {
	GetCourierServicesByCourierIDs(courierIDs []int) ([]*model.CourierService, error)
}

type courierServiceServiceImpl struct {
	courierServiceRepository repository.CourierServiceRepository
}

type CourierServiceSConfig struct {
	CourierServiceRepository repository.CourierServiceRepository
}

func NewCourierServiceService(cfg *CourierServiceSConfig) CourierServiceService {
	return &courierServiceServiceImpl{
		courierServiceRepository: cfg.CourierServiceRepository,
	}
}

func (s *courierServiceServiceImpl) GetCourierServicesByCourierIDs(courierIDs []int) ([]*model.CourierService, error) {
	return s.courierServiceRepository.GetByCourierIDs(courierIDs)
}
