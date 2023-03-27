package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSellerPromotions(t *testing.T) {
	type input struct {
		userID  int
		request *dto.SellerPromotionFilterRequest
	}
	type expected struct {
		data *commonDto.PaginationResponse
		err  error
	}

	var (
		userID     = 1
		shopID     = 1
		limit      = 20
		page       = 1
		request    = &dto.SellerPromotionFilterRequest{Limit: limit, Page: page}
		promotions = []*dto.SellerPromotion{}
		totalRows  = int64(0)
		totalPages = 0
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ShopPromotionRepository)
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get promotions",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{UserID: userID, ID: shopID}, nil)
				pr.On("GetSellerPromotions", shopID, request).Return(nil, int64(0), 0, errors.New("failed to get promotions"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get promotions"),
			},
		},
		{
			description: "should return promotions data when succeed to get promotions",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{UserID: userID, ID: shopID}, nil)
				pr.On("GetSellerPromotions", shopID, request).Return(promotions, totalRows, totalPages, nil)
			},
			expected: expected{
				data: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       promotions,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			shopPromotionRepository := mocks.NewShopPromotionRepository(t)
			tc.beforeTest(shopService, shopPromotionRepository)
			shopPromotionService := service.NewShopPromotionService(&service.ShopPromotionSConfig{
				ShopService:             shopService,
				ShopPromotionRepository: shopPromotionRepository,
			})

			data, err := shopPromotionService.GetSellerPromotions(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetSellerPromotionById(t *testing.T) {
	type input struct {
		userID      int
		promotionId int
	}
	type expected struct {
		data *dto.SellerPromotion
		err  error
	}

	var (
		userID      = 1
		shopID      = 1
		promotionId = 1
		promotion   = dto.SellerPromotion{}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ShopPromotionRepository)
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				userID:      userID,
				promotionId: promotionId,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get promotion",
			input: input{
				userID:      userID,
				promotionId: promotionId,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID, UserID: userID}, nil)
				vr.On("GetSellerPromotionById", shopID, promotionId).Return(nil, errors.New("failed to get promotion"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get promotion"),
			},
		},
		{
			description: "should return dto when succeed to get promotion",
			input: input{
				userID:      userID,
				promotionId: promotionId,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID, UserID: userID}, nil)
				vr.On("GetSellerPromotionById", shopID, promotionId).Return(&promotion, nil)
			},
			expected: expected{
				data: &dto.SellerPromotion{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			shopPromotionRepository := mocks.NewShopPromotionRepository(t)
			tc.beforeTest(shopService, shopPromotionRepository)
			shopPromotionService := service.NewShopPromotionService(&service.ShopPromotionSConfig{
				ShopService:             shopService,
				ShopPromotionRepository: shopPromotionRepository,
			})

			data, err := shopPromotionService.GetSellerPromotionById(tc.input.userID, tc.input.promotionId)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}
