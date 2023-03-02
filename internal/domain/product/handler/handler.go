package handler

import "kedai/backend/be-kedai/internal/domain/product/service"

type Handler struct {
	categoryService service.CategoryService
	productService  service.ProductService
	skuSerivce      service.SkuService
}

type Config struct {
	CategoryService service.CategoryService
	ProductService  service.ProductService
	SkuSerivce      service.SkuService
}

func New(cfg *Config) *Handler {
	return &Handler{
		categoryService: cfg.CategoryService,
		productService:  cfg.ProductService,
		skuSerivce:      cfg.SkuSerivce,
	}
}
