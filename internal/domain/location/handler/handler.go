package handler

import "kedai/backend/be-kedai/internal/domain/location/service"

type Handler struct {
	cityService     service.CityService
	provinceService service.ProvinceService
}

type Config struct {
	CityService     service.CityService
	ProvinceService service.ProvinceService
}

func New(cfg *Config) *Handler {
	return &Handler{
		cityService:     cfg.CityService,
		provinceService: cfg.ProvinceService,
	}
}
