package handler

import (
	locationService "kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/internal/domain/user/service"
)

type Handler struct {
	userService          service.UserService
	userWishlistService  service.UserWishlistService
	userCartItemService  service.UserCartItemService
	walletService        service.WalletService
	walletHistoryService service.WalletHistoryService
	sealabsPayService    service.SealabsPayService
	addressService       locationService.AddressService
	userProfileService   service.UserProfileService
}

type HandlerConfig struct {
	UserService          service.UserService
	UserWishlistService  service.UserWishlistService
	UserCartItemService  service.UserCartItemService
	WalletService        service.WalletService
	WalletHistoryService service.WalletHistoryService
	SealabsPayService    service.SealabsPayService
	AddressService       locationService.AddressService
	UserProfileService   service.UserProfileService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:          cfg.UserService,
		userWishlistService:  cfg.UserWishlistService,
		userCartItemService:  cfg.UserCartItemService,
		walletService:        cfg.WalletService,
		walletHistoryService: cfg.WalletHistoryService,
		sealabsPayService:    cfg.SealabsPayService,
		addressService:       cfg.AddressService,
		userProfileService:   cfg.UserProfileService,
	}
}
