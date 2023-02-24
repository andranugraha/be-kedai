package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userService         service.UserService
	userAddressService  service.UserAddressService
	userWishlistService service.UserWishlistService
	userCartItemService service.UserCartItemService
	walletService       service.WalletService
}

type HandlerConfig struct {
	UserService         service.UserService
	UserAddressService  service.UserAddressService
	UserWishlistService service.UserWishlistService
	UserCartItemService service.UserCartItemService
	WalletService       service.WalletService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:         cfg.UserService,
		userAddressService:  cfg.UserAddressService,
		userWishlistService: cfg.UserWishlistService,
		userCartItemService: cfg.UserCartItemService,
		walletService:       cfg.WalletService,
	}
}
