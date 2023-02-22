package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userService   service.UserService
	walletService service.WalletService
}

type HandlerConfig struct {
	UserService   service.UserService
	WalletService service.WalletService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:   cfg.UserService,
		walletService: cfg.WalletService,
	}
}
