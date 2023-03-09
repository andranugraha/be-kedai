package service

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/repository"
)

type AddressService interface {
	SearchAddress(req *dto.SearchAddressRequest) ([]*dto.SearchAddressResponse, error)
}

type addressServiceImpl struct {
	addressRepo        repository.AddressRepository
	subdistrictService SubdistrictService
}

type AddressSConfig struct {
	AddressRepo        repository.AddressRepository
	SubdistrictService SubdistrictService
}

func NewAddressService(cfg *AddressSConfig) AddressService {
	return &addressServiceImpl{
		addressRepo:        cfg.AddressRepo,
		subdistrictService: cfg.SubdistrictService,
	}
}

func (s *addressServiceImpl) SearchAddress(req *dto.SearchAddressRequest) ([]*dto.SearchAddressResponse, error) {
	return s.addressRepo.SearchAddress(req)
}
