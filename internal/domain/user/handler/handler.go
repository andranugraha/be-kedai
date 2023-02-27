package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userService         service.UserService
	userWishlistService service.UserWishlistService
	userCartItemService service.UserCartItemService
	walletService       service.WalletService
	sealabsPayService   service.SealabsPayService
}

type HandlerConfig struct {
	UserService         service.UserService
	UserWishlistService service.UserWishlistService
	UserCartItemService service.UserCartItemService
	WalletService       service.WalletService
	SealabsPayService   service.SealabsPayService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:         cfg.UserService,
		userWishlistService: cfg.UserWishlistService,
		userCartItemService: cfg.UserCartItemService,
		walletService:       cfg.WalletService,
		sealabsPayService:   cfg.SealabsPayService,
	}
}
