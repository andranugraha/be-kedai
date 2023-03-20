package service

import (
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/repository"
)

type ChatService interface {
	GetAllRoom(userId int, page int, limit int) []*dto.ChatListResponse
	GetAllChat(roomId string, page int, limit int) []*dto.ChatDetailResponse
	AddChat(roomId string, msg string, chatType string, userId int) (*dto.ChatDetailResponse, error)
	ForceReadRoom(roomId string) *dto.ChatListResponse
}

type chatServiceImpl struct {
	chatRepo repository.ChatRepository
}

type ChatConfig struct {
	ChatRepo repository.ChatRepository
}

func NewChatService(config *ChatConfig) ChatService {
	return &chatServiceImpl{
		chatRepo: config.ChatRepo,
	}
}

func (s *chatServiceImpl) GetAllRoom(userId int, page int, limit int) []*dto.ChatListResponse {
	return []*dto.ChatListResponse{}
}

func (s *chatServiceImpl) GetAllChat(roomId string, page int, limit int) []*dto.ChatDetailResponse {
	return []*dto.ChatDetailResponse{}
}

func (s *chatServiceImpl) AddChat(roomId string, msg string, chatType string, userId int) (*dto.ChatDetailResponse, error) {
	return s.chatRepo.AddChat(roomId, msg, chatType, userId)
}

func (s *chatServiceImpl) ForceReadRoom(roomId string) *dto.ChatListResponse {
	return nil
}
