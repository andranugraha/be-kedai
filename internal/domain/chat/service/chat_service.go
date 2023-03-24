package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/repository"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
)

type ChatService interface {
	UserGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.UserListOfChatResponse, error)
	SellerGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.SellerListOfChatResponse, error)
	UserGetChat(param *dto.ChatParamRequest, userId int, shopSlug string) (*commonDto.PaginationResponse, error)
	SellerGetChat(param *dto.ChatParamRequest, userId int, username string) (*commonDto.PaginationResponse, error)
	UserAddChat(body *dto.SendChatBodyRequest, userId int, shopSlug string) (*dto.ChatResponse, error)
	SellerAddChat(body *dto.SendChatBodyRequest, userId int, username string) (*dto.ChatResponse, error)
}

type chatServiceImpl struct {
	chatRepo    repository.ChatRepository
	shopService shopService.ShopService
	userService userService.UserService
}

type ChatConfig struct {
	ChatRepo    repository.ChatRepository
	ShopService shopService.ShopService
	UserService userService.UserService
}

func NewChatService(config *ChatConfig) ChatService {
	return &chatServiceImpl{
		chatRepo:    config.ChatRepo,
		shopService: config.ShopService,
		userService: config.UserService,
	}
}

func (s *chatServiceImpl) UserGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.UserListOfChatResponse, error) {
	return s.chatRepo.UserGetListOfChats(param, userId)
}

func (s *chatServiceImpl) SellerGetListOfChats(param *dto.ListOfChatsParamRequest, userId int) ([]*dto.SellerListOfChatResponse, error) {
	return s.chatRepo.SellerGetListOfChats(param, userId)
}

func (s *chatServiceImpl) UserGetChat(param *dto.ChatParamRequest, userId int, shopSlug string) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopBySlug(shopSlug)
	if err != nil {
		return nil, err
	}
	return s.chatRepo.UserGetChat(param, userId, shop)
}

func (s *chatServiceImpl) SellerGetChat(param *dto.ChatParamRequest, userId int, username string) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.chatRepo.SellerGetChat(param, shop, user)
}

func (s *chatServiceImpl) UserAddChat(body *dto.SendChatBodyRequest, userId int, shopSlug string) (*dto.ChatResponse, error) {
	shop, err := s.shopService.FindShopBySlug(shopSlug)
	if err != nil {
		return nil, err
	}

	if shop.UserID == userId {
		return nil, errs.ErrSelfMessaging
	}

	if body.Type == "product" {
		// Check Message as ProductCode: ErrProductNotFound
	}

	if body.Type == "invoice" || body.Type == "complaint" {
		// Check Message as InvoiceCode: ErrInvoiceNotFound
	}

	return s.chatRepo.UserAddChat(body, userId, shop)
}

func (s *chatServiceImpl) SellerAddChat(body *dto.SendChatBodyRequest, userId int, username string) (*dto.ChatResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if shop.UserID == user.ID {
		return nil, errs.ErrSelfMessaging
	}

	if body.Type == "product" {
		// Check Message as ProductCode: ErrProductNotFound
	}

	if body.Type == "invoice" || body.Type == "complaint" {
		// Check Message as InvoiceCode: ErrInvoiceNotFound
	}

	return s.chatRepo.SellerAddChat(body, shop, user)
}
