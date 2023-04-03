package handler

import "kedai/backend/be-kedai/internal/domain/marketplace/service"

type Handler struct {
	marketplaceVoucherService service.MarketplaceVoucherService
	marketplaceBannerService  service.MarketplaceBannerService
}

type HandlerConfig struct {
	MarketplaceVoucherService service.MarketplaceVoucherService
	MarketplaceBannerService  service.MarketplaceBannerService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		marketplaceVoucherService: cfg.MarketplaceVoucherService,
		marketplaceBannerService:  cfg.MarketplaceBannerService,
	}
}
