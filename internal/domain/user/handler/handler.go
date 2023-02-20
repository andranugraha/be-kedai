package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userWishlistService service.UserWishlistService
}

type HandlerConfig struct {
	UserWishlistService service.UserWishlistService
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		userWishlistService: cfg.UserWishlistService,
	}
}
