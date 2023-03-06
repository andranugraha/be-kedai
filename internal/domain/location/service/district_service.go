package service

import (
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
}

type DistrictSConfig struct {
	DistrictRepo repository.DistrictRepository
}

func NewDistrictService(cfg *DistrictSConfig) DistrictService {
	return &districtServiceImpl{
		districtRepo: cfg.DistrictRepo,
	}
}

func (d *districtServiceImpl) GetDistrictByID(districtID int) (district *model.District, err error) {
	return d.districtRepo.GetByID(districtID)
}

func (d *districtServiceImpl) GetDistricts(req dto.GetDistrictsRequest) (districts []*model.District, err error) {
	return d.districtRepo.GetAll(req)
}
