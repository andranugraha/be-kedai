package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/service"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserGetChat(t *testing.T) {
	type input struct {
		param    *dto.ChatParamRequest
		userId   int
		shopSlug string
		mockData *commonDto.PaginationResponse
		mockErr  error
	}
	type expected struct {
		shop *model.Shop
		data *commonDto.PaginationResponse
		err  error
	}

	var (
		param = &dto.ChatParamRequest{
			LimitByDay: 366,
			Page:       1,
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
				param:    param,
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
				param:    param,
				userId:   userId,
				shopSlug: shopSlug,
				mockData: &commonDto.PaginationResponse{Data: &dto.ChatResponse{}},
				mockErr:  nil,
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(shop, nil)
				cr.On("UserGetChat", param, userId, shop).Return(&commonDto.PaginationResponse{Data: &dto.ChatResponse{}}, nil)
			},
			expected: expected{
				data: &commonDto.PaginationResponse{Data: &dto.ChatResponse{}},
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

		data, err := chatService.UserGetChat(tc.input.param, tc.input.userId, tc.input.shopSlug)

		assert.Equal(t, tc.expected.data, data)
		assert.Equal(t, tc.expected.err, err)
	}
}

func TestSellerGetChat(t *testing.T) {
	type input struct {
		param    *dto.ChatParamRequest
		userId   int
		username string
		mockData *commonDto.PaginationResponse
		mockErr  error
	}
	type expected struct {
		shop *model.Shop
		data *commonDto.PaginationResponse
		err  error
	}

	var (
		param = &dto.ChatParamRequest{
			LimitByDay: 366,
			Page:       1,
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
				param:    param,
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
				param:    param,
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
				param:    param,
				userId:   userId,
				username: username,
				mockData: &commonDto.PaginationResponse{Data: &dto.ChatResponse{}},
				mockErr:  nil,
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService) {
				ss.On("FindShopByUserId", userId).Return(shop, nil)
				us.On("GetByUsername", username).Return(user, nil)
				cr.On("SellerGetChat", param, shop, user).Return(&commonDto.PaginationResponse{Data: &dto.ChatResponse{}}, nil)
			},
			expected: expected{
				data: &commonDto.PaginationResponse{Data: &dto.ChatResponse{}},
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

		data, err := chatService.SellerGetChat(tc.input.param, tc.input.userId, tc.input.username)

		assert.Equal(t, tc.expected.data, data)
		assert.Equal(t, tc.expected.err, err)
	}
}

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
			Type:    "text",
		}
		userId   = 1
		shopSlug = "shop-A"
		shop     = &model.Shop{UserID: 2}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ChatRepository, *mocks.ShopService, *mocks.ProductService, *mocks.InvoicePerShopService)
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
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(shop, errors.New("failed to find shop"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("failed to find shop"),
			},
		},
		{
			description: "should return error when self messaging",
			input: input{
				body:     body,
				userId:   userId,
				shopSlug: shopSlug,
				mockData: nil,
				mockErr:  errs.ErrSelfMessaging,
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(&model.Shop{UserID: 1}, nil)
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errs.ErrSelfMessaging,
			},
		},
		{
			description: "should return error when type of message is product and product not found",
			input: input{
				body: &dto.SendChatBodyRequest{
					Type:    "product",
					Message: "ITEM-001",
				},
				userId:   userId,
				shopSlug: shopSlug,
				mockData: nil,
				mockErr:  errors.New("product not found"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(shop, nil)
				ps.On("GetByCode", "ITEM-001").Return(nil, errors.New("product not found"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("product not found"),
			},
		},
		{
			description: "should return error when type of message is invoice and invoice not found",
			input: input{
				body: &dto.SendChatBodyRequest{
					Type:    "invoice",
					Message: "INV-A",
				},
				userId:   userId,
				shopSlug: shopSlug,
				mockData: nil,
				mockErr:  errors.New("invoice not found"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopBySlug", shopSlug).Return(shop, nil)
				is.On("GetInvoicesByUserIDAndCode", userId, "INV-A").Return(nil, errors.New("invoice not found"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("invoice not found"),
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
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
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
		mockProductService := mocks.NewProductService(t)
		mockInvoiceService := mocks.NewInvoicePerShopService(t)
		mockChatRepo := mocks.NewChatRepository(t)
		tc.beforeTest(mockChatRepo, mockShopService, mockProductService, mockInvoiceService)
		chatService := service.NewChatService(&service.ChatConfig{
			ChatRepo:       mockChatRepo,
			ShopService:    mockShopService,
			ProductService: mockProductService,
			InvoiceService: mockInvoiceService,
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
			Type:    "text",
		}
		userId   = 1
		username = "usernameA"
		shop     = &model.Shop{UserID: 1}
		user     = &userModel.User{ID: 2}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ChatRepository, *mocks.ShopService, *mocks.UserService, *mocks.ProductService, *mocks.InvoicePerShopService)
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
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
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
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
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
			description: "should return error when self messaging",
			input: input{
				body:     body,
				userId:   userId,
				username: username,
				mockData: nil,
				mockErr:  errs.ErrSelfMessaging,
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopByUserId", userId).Return(shop, nil)
				us.On("GetByUsername", username).Return(&userModel.User{ID: 1}, nil)
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errs.ErrSelfMessaging,
			},
		},
		{
			description: "should return error when type of message is product and product not found",
			input: input{
				body: &dto.SendChatBodyRequest{
					Type:    "product",
					Message: "ITEM-001",
				},
				userId:   userId,
				username: username,
				mockData: nil,
				mockErr:  errors.New("product not found"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopByUserId", userId).Return(shop, nil)
				us.On("GetByUsername", username).Return(user, nil)
				ps.On("GetByCode", "ITEM-001").Return(nil, errors.New("product not found"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("product not found"),
			},
		},
		{
			description: "should return error when type of message is invoice and invoice not found",
			input: input{
				body: &dto.SendChatBodyRequest{
					Type:    "invoice",
					Message: "INV-A",
				},
				userId:   userId,
				username: username,
				mockData: nil,
				mockErr:  errors.New("invoice not found"),
			},
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
				ss.On("FindShopByUserId", userId).Return(shop, nil)
				us.On("GetByUsername", username).Return(user, nil)
				is.On("GetInvoicesByUserIDAndCode", user.ID, "INV-A").Return(nil, errors.New("invoice not found"))
			},
			expected: expected{
				shop: shop,
				data: nil,
				err:  errors.New("invoice not found"),
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
			beforeTest: func(cr *mocks.ChatRepository, ss *mocks.ShopService, us *mocks.UserService, ps *mocks.ProductService, is *mocks.InvoicePerShopService) {
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
		mockProductService := mocks.NewProductService(t)
		mockInvoiceService := mocks.NewInvoicePerShopService(t)
		mockChatRepo := mocks.NewChatRepository(t)
		tc.beforeTest(mockChatRepo, mockShopService, mockUserService, mockProductService, mockInvoiceService)
		chatService := service.NewChatService(&service.ChatConfig{
			ChatRepo:       mockChatRepo,
			ShopService:    mockShopService,
			UserService:    mockUserService,
			ProductService: mockProductService,
			InvoiceService: mockInvoiceService,
		})

		data, err := chatService.SellerAddChat(tc.input.body, tc.input.userId, tc.input.username)

		assert.Equal(t, tc.expected.data, data)
		assert.Equal(t, tc.expected.err, err)
	}
}
