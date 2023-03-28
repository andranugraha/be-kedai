package handler

import (
	orderService "kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/internal/domain/product/service"
)

type Handler struct {
	categoryService          service.CategoryService
	productService           service.ProductService
	skuService               service.SkuService
	transactionReviewService orderService.TransactionReviewService
	discussionService        service.DiscussionService
}

type Config struct {
	CategoryService          service.CategoryService
	ProductService           service.ProductService
	SkuService               service.SkuService
	TransactionReviewService orderService.TransactionReviewService
	DiscussionService        service.DiscussionService
}

func New(cfg *Config) *Handler {
	return &Handler{
		categoryService:          cfg.CategoryService,
		productService:           cfg.ProductService,
		skuService:               cfg.SkuService,
		transactionReviewService: cfg.TransactionReviewService,
		discussionService:        cfg.DiscussionService,
	}
}