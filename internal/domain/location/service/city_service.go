package service

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type CityService interface {
	GetCities(dto.GetCitiesRequest) (*dto.GetCitiesResponse, error)
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

func (c *cityServiceImpl) GetCities(req dto.GetCitiesRequest) (res *dto.GetCitiesResponse, err error) {
	cities, totalRows, totalPages, err := c.cityRepo.GetAll(req)
	if err != nil {
		return
	}

	res = &dto.GetCitiesResponse{
		Data:       cities,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Limit:      req.Limit,
		Page:       req.Page,
	}

	return
}
