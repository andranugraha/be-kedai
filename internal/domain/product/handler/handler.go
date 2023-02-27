package handler

import "kedai/backend/be-kedai/internal/domain/product/service"

type Handler struct {
	categoryService service.CategoryService
	producteService service.ProductService
}

type Config struct {
	CategoryService service.CategoryService
	ProductService  service.ProductService
}

func New(cfg *Config) *Handler {
	return &Handler{
		categoryService: cfg.CategoryService,
		producteService: cfg.ProductService,
	}
}
