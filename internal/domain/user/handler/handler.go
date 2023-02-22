package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
<<<<<<< HEAD
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
=======
	userService service.UserService
}

type HandlerConfig struct {
	UserService service.UserService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService: cfg.UserService,
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
	}
}
