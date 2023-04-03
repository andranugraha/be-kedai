package handler

import "kedai/backend/be-kedai/internal/domain/shop/service"

type Handler struct {
	shopService          service.ShopService
	shopVoucherService   service.ShopVoucherService
	shopPromotionService service.ShopPromotionService
	courierService       service.CourierService
	shopGuestService     service.ShopGuestService
	shopCategoryService  service.ShopCategoryService
}

type HandlerConfig struct {
	ShopService          service.ShopService
	ShopVoucherService   service.ShopVoucherService
	ShopPromotionService service.ShopPromotionService
	CourierService       service.CourierService
	ShopGuestService     service.ShopGuestService
	ShopCategoryService  service.ShopCategoryService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		shopService:          cfg.ShopService,
		shopVoucherService:   cfg.ShopVoucherService,
		shopPromotionService: cfg.ShopPromotionService,
		courierService:       cfg.CourierService,
		shopGuestService:     cfg.ShopGuestService,
		shopCategoryService:  cfg.ShopCategoryService,
	}
}
