package handler

import "kedai/backend/be-kedai/internal/domain/product/service"

type Handler struct {
	categoryService service.CategoryService
}

type Config struct {
	CategoryService service.CategoryService
}

func New(cfg *Config) *Handler {
	return &Handler{
		categoryService: cfg.CategoryService,
	}
}
