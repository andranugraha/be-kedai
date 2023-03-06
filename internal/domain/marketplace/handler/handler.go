package handler

import "kedai/backend/be-kedai/internal/domain/marketplace/service"

type Handler struct {
	marketplaceVoucherService service.MarketplaceVoucherService
}

type HandlerConfig struct {
	MarketplaceVoucherService service.MarketplaceVoucherService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		marketplaceVoucherService: cfg.MarketplaceVoucherService,
	}
}
