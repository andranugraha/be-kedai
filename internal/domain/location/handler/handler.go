package handler

import "kedai/backend/be-kedai/internal/domain/location/service"

type Handler struct {
	cityService        service.CityService
	provinceService    service.ProvinceService
	districtService    service.DistrictService
	subdistrictService service.SubdistrictService
	addressService     service.AddressService
}

type Config struct {
	CityService        service.CityService
	ProvinceService    service.ProvinceService
	DistrictService    service.DistrictService
	SubdistrictService service.SubdistrictService
	AddressService     service.AddressService
}

func New(cfg *Config) *Handler {
	return &Handler{
		cityService:        cfg.CityService,
		provinceService:    cfg.ProvinceService,
		districtService:    cfg.DistrictService,
		subdistrictService: cfg.SubdistrictService,
		addressService:     cfg.AddressService,
	}
}
