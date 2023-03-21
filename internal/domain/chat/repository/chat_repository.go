package repository

import (
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
	UserGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.UserListOfChatResponse, error)
	SellerGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.SellerListOfChatResponse, error)
	UserGetChat(param *dto.ChatParamRequest, userId int, shopSlug string) ([]*dto.ChatResponse, error)
	SellerGetChat(param *dto.ChatParamRequest, userId int, username string) ([]*dto.ChatResponse, error)
	UserAddChat(body *dto.SendChatBodyRequest, userId int, shop *shopModel.Shop) (*dto.ChatResponse, error)
	SellerAddChat(body *dto.SendChatBodyRequest, shop *shopModel.Shop, user *userModel.User) (*dto.ChatResponse, error)
}

type chatRepositoryImpl struct {
	db *gorm.DB
}

type ChatRConfig struct {
	DB *gorm.DB
}

func NewChatRepository(cfg *ChatRConfig) ChatRepository {
	return &chatRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *chatRepositoryImpl) Last(chat *model.Chat) (*model.Chat, error) {
	result := r.db.Last(&chat)
	return chat, result.Error
}

func (r *chatRepositoryImpl) UserGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.UserListOfChatResponse, error) {
	return []*dto.UserListOfChatResponse{}, nil
}

func (r *chatRepositoryImpl) SellerGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.SellerListOfChatResponse, error) {
	return []*dto.SellerListOfChatResponse{}, nil
}

func (r *chatRepositoryImpl) UserGetChat(param *dto.ChatParamRequest, userId int, shopSlug string) ([]*dto.ChatResponse, error) {
	return []*dto.ChatResponse{}, nil
}

func (r *chatRepositoryImpl) SellerGetChat(param *dto.ChatParamRequest, userId int, username string) ([]*dto.ChatResponse, error) {
	return []*dto.ChatResponse{}, nil
}

func (r *chatRepositoryImpl) UserAddChat(body *dto.SendChatBodyRequest, userId int, shop *shopModel.Shop) (*dto.ChatResponse, error) {
	if body.Type == "" {
		body.Type = "normal"
	}

	chat := &model.Chat{
		Message: body.Message,
		Type:    body.Type,
		ShopId:  shop.ID,
		UserId:  userId,
		Issuer:  "user",
	}
	result := r.db.Create(&chat)
	if result.Error != nil {
		return nil, result.Error
	}
	chat, _ = r.Last(chat)
	return dto.ConvertChatToOutput(chat, "user"), nil
}

func (r *chatRepositoryImpl) SellerAddChat(body *dto.SendChatBodyRequest, shop *shopModel.Shop, user *userModel.User) (*dto.ChatResponse, error) {
	if body.Type == "" {
		body.Type = "normal"
	}

	chat := &model.Chat{
		Message: body.Message,
		Type:    body.Type,
		ShopId:  shop.ID,
		UserId:  user.ID,
		Issuer:  "seller",
	}
	result := r.db.Create(&chat)
	if result.Error != nil {
		return nil, result.Error
	}
	chat, _ = r.Last(chat)
	return dto.ConvertChatToOutput(chat, "seller"), nil
}
