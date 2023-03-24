package repository

import (
	"fmt"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/date"
	"math"
	"time"

	"gorm.io/gorm"
)

type ChatRepository interface {
	UserGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.UserListOfChatResponse, error)
	SellerGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.SellerListOfChatResponse, error)
	UserGetChat(param *dto.ChatParamRequest, userId int, shop *shopModel.Shop) (*commonDto.PaginationResponse, error)
	SellerGetChat(param *dto.ChatParamRequest, shop *shopModel.Shop, user *userModel.User) (*commonDto.PaginationResponse, error)
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

func (r *chatRepositoryImpl) Find(shop *shopModel.Shop, user *userModel.User) ([]*model.Chat, error) {
	var chats []*model.Chat
	result := r.db.Where("user_id = ?", user.ID).Where("shop_id = ?", shop.ID).Order("created_at DESC").Find(&chats)
	return chats, result.Error
}

func (r *chatRepositoryImpl) FindWithParam(param *dto.ChatParamRequest, shop *shopModel.Shop, user *userModel.User) ([]*model.Chat, error) {
	offset := param.Offset()
	endDate := time.Now().AddDate(0, 0, -1*offset)
	startDate := endDate.AddDate(0, 0, -1*param.LimitByDay)
	var chats []*model.Chat
	result := r.db.Where("user_id = ?", user.ID).Where("shop_id = ?", shop.ID).Where("created_at >= ? AND created_at <= ?", startDate, endDate).Order("created_at DESC").Find(&chats)
	return chats, result.Error
}

func (r *chatRepositoryImpl) FirstChat(shop *shopModel.Shop, user *userModel.User) (*model.Chat, error) {
	var chat *model.Chat
	result := r.db.Where("user_id = ?", user.ID).Where("shop_id = ?", shop.ID).Order("created_at ASC").First(&chat)
	return chat, result.Error
}

func (r *chatRepositoryImpl) UserGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.UserListOfChatResponse, error) {
	return []*dto.UserListOfChatResponse{}, nil
}

func (r *chatRepositoryImpl) SellerGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.SellerListOfChatResponse, error) {
	return []*dto.SellerListOfChatResponse{}, nil
}

func (r *chatRepositoryImpl) UserGetChat(param *dto.ChatParamRequest, userId int, shop *shopModel.Shop) (*commonDto.PaginationResponse, error) {
	var calculatedTotalRows int64
	firstChat, err := r.FirstChat(shop, &userModel.User{ID: userId})
	if err != nil && err == gorm.ErrRecordNotFound {
		calculatedTotalRows = 0
	} else {
		calculatedTotalRows = int64(date.DaysBetween(firstChat.CreatedAt, time.Now()))
		fmt.Println("aASD", date.DaysBetween(firstChat.CreatedAt, time.Now()))
		fmt.Println("1", firstChat.CreatedAt)
		fmt.Println("2", time.Now())
	}

	var chats []*model.Chat
	chats, err = r.FindWithParam(param, shop, &userModel.User{ID: userId})
	if err != nil {
		return nil, err
	}

	chatsResponse := []*dto.ChatResponse{}
	for _, chat := range chats {
		chatsResponse = append(chatsResponse, dto.ConvertChatToOutput(chat, "user"))
	}

	paginatedChats := &commonDto.PaginationResponse{
		Data:       chatsResponse,
		Page:       param.Page,
		Limit:      param.LimitByDay,
		TotalRows:  calculatedTotalRows,
		TotalPages: int(math.Ceil(float64(calculatedTotalRows) / float64(param.LimitByDay))),
	}
	return paginatedChats, nil
}

func (r *chatRepositoryImpl) SellerGetChat(param *dto.ChatParamRequest, shop *shopModel.Shop, user *userModel.User) (*commonDto.PaginationResponse, error) {
	var calculatedTotalRows int64
	firstChat, err := r.FirstChat(shop, user)
	if err != nil && err == gorm.ErrRecordNotFound {
		calculatedTotalRows = 0
	} else {
		calculatedTotalRows = int64(date.DaysBetween(firstChat.CreatedAt, time.Now()))
		fmt.Println("aASD", date.DaysBetween(firstChat.CreatedAt, time.Now()))
		fmt.Println("1", firstChat.CreatedAt)
		fmt.Println("2", time.Now())
	}

	var chats []*model.Chat
	chats, err = r.FindWithParam(param, shop, user)
	if err != nil {
		return nil, err
	}

	chatsResponse := []*dto.ChatResponse{}
	for _, chat := range chats {
		chatsResponse = append(chatsResponse, dto.ConvertChatToOutput(chat, "seller"))
	}

	paginatedChats := &commonDto.PaginationResponse{
		Data:       chatsResponse,
		Page:       param.Page,
		Limit:      param.LimitByDay,
		TotalRows:  calculatedTotalRows,
		TotalPages: int(math.Ceil(float64(calculatedTotalRows) / float64(param.LimitByDay))),
	}
	return paginatedChats, nil
}

func (r *chatRepositoryImpl) UserAddChat(body *dto.SendChatBodyRequest, userId int, shop *shopModel.Shop) (*dto.ChatResponse, error) {
	if body.Type == "" {
		body.Type = "text"
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
		body.Type = "text"
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
