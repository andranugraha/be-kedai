package handler

import "kedai/backend/be-kedai/internal/domain/shop/service"

type Handler struct {
	shopService service.ShopService
}

type HandlerConfig struct {
	ShopService service.ShopService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		shopService: cfg.ShopService,
	}
}