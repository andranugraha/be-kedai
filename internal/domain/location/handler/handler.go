package handler

import "kedai/backend/be-kedai/internal/domain/location/service"

type Handler struct {
	cityService service.CityService
}

type Config struct {
	CityService service.CityService
}

func New(cfg *Config) *Handler {
	return &Handler{
		cityService: cfg.CityService,
	}
}
