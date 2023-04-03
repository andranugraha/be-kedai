package service

import (
	"kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/location/cache"
	locationDto "kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type CityService interface {
	GetCities(locationDto.GetCitiesRequest) (*dto.PaginationResponse, error)
	GetCityByID(int) (*model.City, error)
}

type cityServiceImpl struct {
	cityRepo repository.CityRepository
	cache    cache.LocationCache
}

type CitySConfig struct {
	CityRepo repository.CityRepository
	Cache    cache.LocationCache
}

func NewCityService(cfg *CitySConfig) CityService {
	return &cityServiceImpl{
		cityRepo: cfg.CityRepo,
		cache:    cfg.Cache,
	}
}

func (c *cityServiceImpl) GetCities(req locationDto.GetCitiesRequest) (res *dto.PaginationResponse, err error) {
	res = c.cache.GetCities(req)
	if res != nil {
		return
	}

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

	c.cache.StoreCities(req, res)

	return
}

func (c *cityServiceImpl) GetCityByID(cityID int) (city *model.City, err error) {
	return c.cityRepo.GetByID(cityID)
}
