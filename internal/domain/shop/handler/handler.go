package handler

import "kedai/backend/be-kedai/internal/domain/shop/service"

type Handler struct {
	shopService        service.ShopService
	shopVoucherService service.ShopVoucherService
}

type HandlerConfig struct {
	ShopService        service.ShopService
	ShopVoucherService service.ShopVoucherService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		shopService:        cfg.ShopService,
		shopVoucherService: cfg.ShopVoucherService,
	}
}
