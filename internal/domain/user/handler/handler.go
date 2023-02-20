package handler

import "kedai/backend/be-kedai/internal/domain/user/service"

type Handler struct {
	userService service.UserService
}

type HandlerConfig struct {
	UserService service.UserService
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService: cfg.UserService,
	}
}
