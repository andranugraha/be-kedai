package service

import (
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type CityService interface {
	GetCities() ([]*model.City, error)
}

type cityServiceImpl struct {
	cityRepo repository.CityRepository
}

type CitySConfig struct {
	CityRepo repository.CityRepository
}

func NewCityService(cfg *CitySConfig) CityService {
	return &cityServiceImpl{
		cityRepo: cfg.CityRepo,
	}
}

func (c *cityServiceImpl) GetCities() ([]*model.City, error) {
	return nil, nil
}
