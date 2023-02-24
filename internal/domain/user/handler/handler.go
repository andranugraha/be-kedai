package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userService         service.UserService
	userWishlistService service.UserWishlistService
	userCartItemService service.UserCartItemService
	walletService       service.WalletService
}

type HandlerConfig struct {
	UserService         service.UserService
	UserWishlistService service.UserWishlistService
	UserCartItemService service.UserCartItemService
	WalletService       service.WalletService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:         cfg.UserService,
		userWishlistService: cfg.UserWishlistService,
		userCartItemService: cfg.UserCartItemService,
		walletService:       cfg.WalletService,
	}
}
