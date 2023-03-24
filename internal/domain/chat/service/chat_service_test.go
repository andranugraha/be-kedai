package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/service"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAddChat(t *testing.T) {
	type input struct {
		body     *dto.SendChatBodyRequest
		userId   int
		shopSlug string
		mockData *dto.ChatResponse
		mockErr  error
	}
	type expected struct {
		shop *model.Shop
		data *dto.ChatResponse
		err  error
	}

	var (
		body = &dto.SendChatBodyRequest{
			Message: "hai sayang 1",
			Type:    "complaint",
		}
		userId   = 1
		shopSlug = "shop-A"
		shop     = &model.Shop{}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ChatRepository, *mocks.ShopService)
		expected
	}{
		{
			description: "should return error when failed to find shop slug",
			input: input{
				body:     body,
				userId:   userId,
				shopSlug: shopSlug,
				mockData: nil,
				mockErr:  errors.New("failed to find shop"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(shop, errors.New("failed to find shop"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("failed to find shop"),
			},
		},
		{
			description: "should return data when succeed to send message",
			input: input{
				body:     body,
				userId:   userId,
				shopSlug: shopSlug,
				mockData: &dto.ChatResponse{},
				mockErr:  nil,
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(shop, nil)
				cr.On("UserAddChat", body, userId, shop).Return(&dto.ChatResponse{}, nil)
			},
			expected: expected{
				data: &dto.ChatResponse{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		mockShopService := mocks.NewShopService(t)
		mockChatRepo := mocks.NewChatRepository(t)
		tc.beforeTest(mockChatRepo, mockShopService)
		chatService := service.NewChatService(&service.ChatConfig{
			ChatRepo:    mockChatRepo,
			ShopService: mockShopService,
		})

		data, err := chatService.UserAddChat(tc.input.body, tc.input.userId, tc.input.shopSlug)

		assert.Equal(t, tc.expected.data, data)
		assert.Equal(t, tc.expected.err, err)
	}
}

func TestSellerAddChat(t *testing.T) {
	type input struct {
		body     *dto.SendChatBodyRequest
		userId   int
		username string
		mockData *dto.ChatResponse
		mockErr  error
	}
	type expected struct {
		shop *model.Shop
		data *dto.ChatResponse
		err  error
	}

	var (
		body = &dto.SendChatBodyRequest{
			Message: "hai sayang 1",
			Type:    "complaint",
		}
		userId   = 1
		username = "usernameA"
		shop     = &model.Shop{}
		user     = &userModel.User{}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ChatRepository, *mocks.ShopService, *mocks.UserService)
		expected
	}{
		{
			description: "should return error when failed to find shop slug",
			input: input{
				body:     body,
				userId:   userId,
				username: username,
				mockData: nil,
				mockErr:  errors.New("failed to find shop"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService) {
				ss.On("FindShopByUserId", userId).Return(shop, errors.New("failed to find shop"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("failed to find shop"),
			},
		},
		{
			description: "should return error when failed to find username",
			input: input{
				body:     body,
				userId:   userId,
				username: username,
				mockData: nil,
				mockErr:  errors.New("failed to find user"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService) {
				ss.On("FindShopByUserId", userId).Return(shop, nil)
				us.On("GetByUsername", username).Return(user, errors.New("failed to find user"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("failed to find user"),
			},
		},
		{
			description: "should return data when succeed to send message",
			input: input{
				body:     body,
				userId:   userId,
				username: username,
				mockData: &dto.ChatResponse{},
				mockErr:  nil,
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService) {
				ss.On("FindShopByUserId", userId).Return(shop, nil)
				us.On("GetByUsername", username).Return(user, nil)
				cr.On("SellerAddChat", body, shop, user).Return(&dto.ChatResponse{}, nil)
			},
			expected: expected{
				data: &dto.ChatResponse{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		mockUserService := mocks.NewUserService(t)
		mockShopService := mocks.NewShopService(t)
		mockChatRepo := mocks.NewChatRepository(t)
		tc.beforeTest(mockChatRepo, mockShopService, mockUserService)
		chatService := service.NewChatService(&service.ChatConfig{
			ChatRepo:    mockChatRepo,
			ShopService: mockShopService,
			UserService: mockUserService,
		})

		data, err := chatService.SellerAddChat(tc.input.body, tc.input.userId, tc.input.username)

		assert.Equal(t, tc.expected.data, data)
		assert.Equal(t, tc.expected.err, err)
	}
}
