package repository

import (
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
	GetAllRoom(userId int, page int, limit int) []*dto.ChatListResponse
	GetAllChat(roomId string, page int, limit int) []*dto.ChatDetailResponse
	AddChat(roomId string, msg string, chatType string, userId int) (*dto.ChatDetailResponse, error)
	ForceReadRoom(roomId string) *dto.ChatListResponse
}

type chatRepositoryImpl struct {
	db *gorm.DB
}

type ChatRConfig struct {
	DB *gorm.DB
}

func NewAddressRepository(cfg *ChatRConfig) ChatRepository {
	return &chatRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *chatRepositoryImpl) Last(chat *model.Chat) (*model.Chat, error) {
	result := r.db.Preload("User.Profile").Last(&chat)
	return chat, result.Error
}

func (r *chatRepositoryImpl) GetAllRoom(userId int, page int, limit int) []*dto.ChatListResponse {
	return []*dto.ChatListResponse{}
}

func (r *chatRepositoryImpl) GetAllChat(roomId string, page int, limit int) []*dto.ChatDetailResponse {
	return []*dto.ChatDetailResponse{}
}

func (r *chatRepositoryImpl) AddChat(roomId string, msg string, chatType string, userId int) (*dto.ChatDetailResponse, error) {
	chat := &model.Chat{
		RoomId:  roomId,
		Message: msg,
		Type:    chatType,
		UserId:  userId,
	}
	result := r.db.Create(&chat)
	if result.Error != nil {
		return nil, result.Error
	}
	chat, _ = r.Last(chat)
	return dto.ConvertChatDetailToOutput(chat), nil
}

func (r *chatRepositoryImpl) ForceReadRoom(roomId string) *dto.ChatListResponse {
	return nil
}
