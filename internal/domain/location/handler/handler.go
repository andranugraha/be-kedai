package handler

import "kedai/backend/be-kedai/internal/domain/location/service"

type Handler struct {
	cityService     service.CityService
	provinceService service.ProvinceService
	districtService service.DistrictService
}

type Config struct {
	CityService     service.CityService
	ProvinceService service.ProvinceService
	DistrictService service.DistrictService
}

func New(cfg *Config) *Handler {
	return &Handler{
		cityService:     cfg.CityService,
		provinceService: cfg.ProvinceService,
		districtService: cfg.DistrictService,
	}
}
