package handler

import "kedai/backend/be-kedai/internal/domain/shop/service"

type Handler struct {
	shopService        service.ShopService
	shopVoucherService service.ShopVoucherService
	courierService     service.CourierService
	shopGuestService   service.ShopGuestService
}

type HandlerConfig struct {
	ShopService        service.ShopService
	ShopVoucherService service.ShopVoucherService
	CourierService     service.CourierService
	ShopGuestService   service.ShopGuestService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		shopService:        cfg.ShopService,
		shopVoucherService: cfg.ShopVoucherService,
		courierService:     cfg.CourierService,
		shopGuestService:   cfg.ShopGuestService,
	}
}
