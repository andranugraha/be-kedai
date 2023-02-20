package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userService         service.UserService
	userWishlistService service.UserWishlistService
}

type HandlerConfig struct {
	UserService         service.UserService
	UserWishlistService service.UserWishlistService
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:         cfg.UserService,
		userWishlistService: cfg.UserWishlistService,
	}
}
