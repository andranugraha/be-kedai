package service

import (
	"kedai/backend/be-kedai/internal/domain/location/cache"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type DistrictService interface {
	GetDistrictByID(int) (*model.District, error)
	GetDistricts(dto.GetDistrictsRequest) (districts []*model.District, err error)
}

type districtServiceImpl struct {
	districtRepo repository.DistrictRepository
	cache        cache.LocationCache
}

type DistrictSConfig struct {
	DistrictRepo repository.DistrictRepository
	Cache        cache.LocationCache
}

func NewDistrictService(cfg *DistrictSConfig) DistrictService {
	return &districtServiceImpl{
		districtRepo: cfg.DistrictRepo,
		cache:        cfg.Cache,
	}
}

func (d *districtServiceImpl) GetDistrictByID(districtID int) (district *model.District, err error) {
	return d.districtRepo.GetByID(districtID)
}

func (d *districtServiceImpl) GetDistricts(req dto.GetDistrictsRequest) (districts []*model.District, err error) {
	districts = d.cache.GetDistricts(req)
	if districts != nil {
		return
	}

	districts, err = d.districtRepo.GetAll(req)
	if err != nil {
		return
	}

	d.cache.StoreDistricts(req, districts)

	return
}
