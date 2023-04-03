package service

import (
	"kedai/backend/be-kedai/internal/domain/location/cache"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type ProvinceService interface {
	GetProvinces() ([]*model.Province, error)
	GetProvinceByID(int) (*model.Province, error)
}

type provinceServiceImpl struct {
	provinceRepo repository.ProvinceRepository
	cache        cache.LocationCache
}

type ProvinceSConfig struct {
	ProvinceRepo repository.ProvinceRepository
	Cache        cache.LocationCache
}

func NewProvinceService(cfg *ProvinceSConfig) ProvinceService {
	return &provinceServiceImpl{
		provinceRepo: cfg.ProvinceRepo,
		cache:        cfg.Cache,
	}
}

func (p *provinceServiceImpl) GetProvinces() (provinces []*model.Province, err error) {
	provinces = p.cache.GetProvinces()
	if provinces != nil {
		return
	}

	provinces, err = p.provinceRepo.GetAll()
	if err != nil {
		return
	}

	p.cache.StoreProvinces(provinces)

	return
}

func (p *provinceServiceImpl) GetProvinceByID(provinceID int) (province *model.Province, err error) {
	return p.provinceRepo.GetByID(provinceID)
}
