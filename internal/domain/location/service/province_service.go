package service

import (
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type ProvinceService interface {
	GetProvinces() ([]*model.Province, error)
	GetProvinceByID(int) (*model.Province, error)
}

type provinceServiceImpl struct {
	provinceRepo repository.ProvinceRepository
}

type ProvinceSConfig struct {
	ProvinceRepo repository.ProvinceRepository
}

func NewProvinceService(cfg *ProvinceSConfig) ProvinceService {
	return &provinceServiceImpl{
		provinceRepo: cfg.ProvinceRepo,
	}
}

func (p *provinceServiceImpl) GetProvinces() (provinces []*model.Province, err error) {
	return p.provinceRepo.GetAll()
}

func (p *provinceServiceImpl) GetProvinceByID(provinceID int) (province *model.Province, err error) {
	return p.provinceRepo.GetByID(provinceID)
}
