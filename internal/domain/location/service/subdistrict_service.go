package service

import (
	"kedai/backend/be-kedai/internal/domain/location/cache"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type SubdistrictService interface {
	GetSubdistrictByID(subdistrictID int) (*model.Subdistrict, error)
	GetSubdistricts(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error)
}

type subdistrictServiceImpl struct {
	subdistrictRepo repository.SubdistrictRepository
	cache           cache.LocationCache
}

type SubdistrictSConfig struct {
	SubdistrictRepo repository.SubdistrictRepository
	Cache           cache.LocationCache
}

func NewSubdistrictService(cfg *SubdistrictSConfig) SubdistrictService {
	return &subdistrictServiceImpl{
		subdistrictRepo: cfg.SubdistrictRepo,
		cache:           cfg.Cache,
	}
}

func (s *subdistrictServiceImpl) GetSubdistrictByID(subdistrictID int) (*model.Subdistrict, error) {
	return s.subdistrictRepo.GetByID(subdistrictID)
}

func (s *subdistrictServiceImpl) GetSubdistricts(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error) {
	subdistricts = s.cache.GetSubdistricts(req)
	if subdistricts != nil {
		return
	}

	subdistricts, err = s.subdistrictRepo.GetAll(req)
	if err != nil {
		return
	}

	s.cache.StoreSubdistricts(req, subdistricts)

	return
}
