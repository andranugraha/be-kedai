package handler

import "kedai/backend/be-kedai/internal/domain/marketplace/service"

type Handler struct {
	marketplaceVoucherService service.MarketplaceVoucherService
	categoryService           service.CategotyService
}

type HandlerConfig struct {
	MarketplaceVoucherService service.MarketplaceVoucherService
	CategoryService           service.CategotyService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		marketplaceVoucherService: cfg.MarketplaceVoucherService,
		categoryService:           cfg.CategoryService,
	}
}
