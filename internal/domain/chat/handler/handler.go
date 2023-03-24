package handler

import (
	"kedai/backend/be-kedai/internal/domain/chat/service"
)

type Handler struct {
	chatService service.ChatService
}

type Config struct {
	ChatService service.ChatService
}

func New(cfg *Config) *Handler {
	return &Handler{
		chatService: cfg.ChatService,
	}
}
