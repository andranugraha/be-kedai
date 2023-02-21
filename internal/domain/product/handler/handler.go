package handler

import "kedai/backend/be-kedai/internal/domain/product/service"

type Handler struct {
	productService service.ProductService
}

type HandlerConfig struct {
	ProductService service.ProductService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		productService: cfg.ProductService,
	}
}
