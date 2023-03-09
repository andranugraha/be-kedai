package service

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type SubdistrictService interface {
	GetSubdistrictByID(subdistrictID int) (*model.Subdistrict, error)
	GetSubdistricts(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error)
	GetDetailSubdistrictByName(subdistrictName string) (*model.Subdistrict, error)
}

type subdistrictServiceImpl struct {
	subdistrictRepo repository.SubdistrictRepository
}

type SubdistrictSConfig struct {
	SubdistrictRepo repository.SubdistrictRepository
}

func NewSubdistrictService(cfg *SubdistrictSConfig) SubdistrictService {
	return &subdistrictServiceImpl{
		subdistrictRepo: cfg.SubdistrictRepo,
	}
}

func (s *subdistrictServiceImpl) GetSubdistrictByID(subdistrictID int) (*model.Subdistrict, error) {
	return s.subdistrictRepo.GetByID(subdistrictID)
}

func (s *subdistrictServiceImpl) GetSubdistricts(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error) {
	return s.subdistrictRepo.GetAll(req)
}

func (s *subdistrictServiceImpl) GetDetailSubdistrictByName(subdistrictName string) (*model.Subdistrict, error) {
	return s.subdistrictRepo.GetDetailByName(subdistrictName)
}
