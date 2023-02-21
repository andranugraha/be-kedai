package service

import (
	"kedai/backend/be-kedai/internal/common/dto"
	locationDto "kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type CityService interface {
	GetCities(locationDto.GetCitiesRequest) (*dto.PaginationResponse, error)
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

func (c *cityServiceImpl) GetCities(req locationDto.GetCitiesRequest) (res *dto.PaginationResponse, err error) {
	cities, totalRows, totalPages, err := c.cityRepo.GetAll(req)
	if err != nil {
		return
	}

	res = &dto.PaginationResponse{
		Data:       cities,
		Limit:      req.Limit,
		Page:       req.Page,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	return
}
