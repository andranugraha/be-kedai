package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	productDto "kedai/backend/be-kedai/internal/domain/product/dto"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestUpdatePromotion(t *testing.T) {
	type input struct {
		userID      int
		shopID      int
		promotionId int
		request     dto.UpdateShopPromotionRequest
	}
	type expected struct {
		err error
	}

	var (
		userID        = 1
		shopID        = 1
		promotionID   = 1
		promotionName = "promotion name"
		isActive      = true
		promotion     = dto.SellerPromotion{
			Product: []*productDto.SellerProductPromotionResponse{
				{
					SKUs: []*productModel.Sku{
						{
							ID: 1,
							Promotion: &productModel.ProductPromotion{
								ID: 1,
							},
						},
					},
				},
			},
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ShopPromotionRepository)
		expected
	}{
		{
			description: "should return error when promotion name is invalid",
			input: input{
				userID: userID,
				request: dto.UpdateShopPromotionRequest{
					Name: "2.2",
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopPromotionRepository) {},
			expected: expected{
				err: errs.ErrInvalidPromotionNamePattern,
			},
		},
		{
			description: "should return error when failed to get shop",
			input: input{
				userID: userID,
				request: dto.UpdateShopPromotionRequest{
					Name: promotionName,
				},
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				err: errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get voucher",
			input: input{
				userID:      userID,
				shopID:      shopID,
				promotionId: promotionID,
				request: dto.UpdateShopPromotionRequest{
					Name: promotionName,
				},
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				pr.On("GetSellerPromotionById", shopID, promotionID).Return(nil, errors.New("failed to get promotion"))
			},
			expected: expected{
				err: errors.New("failed to get promotion"),
			},
		},
		{
			description: "should return error when failed to update promotion",
			input: input{
				userID:      userID,
				shopID:      shopID,
				promotionId: promotionID,
				request: dto.UpdateShopPromotionRequest{
					Name:              promotionName,
					ProductPromotions: []*productDto.UpdateProductPromotionRequest{},
				},
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				pr.On("GetSellerPromotionById", shopID, promotionID).Return(&promotion, nil)
				pr.On("Update", &model.ShopPromotion{Name: promotionName, ShopId: shopID}, mock.Anything).Return(errors.New("failed to update promotion"))
			},
			expected: expected{
				err: errors.New("failed to update promotion"),
			},
		},
		{
			description: "should return success when succeed to update promotion",
			input: input{
				userID:      userID,
				shopID:      shopID,
				promotionId: promotionID,
				request: dto.UpdateShopPromotionRequest{
					Name: promotionName,
					ProductPromotions: []*productDto.UpdateProductPromotionRequest{
						{
							Type:          "discount",
							Amount:        1,
							Stock:         10,
							IsActive:      &isActive,
							PurchaseLimit: 1,
							SkuId:         1,
						},
					},
				},
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopPromotionRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				pr.On("GetSellerPromotionById", shopID, promotionID).Return(&promotion, nil)
				pr.On("Update", &model.ShopPromotion{Name: promotionName, ShopId: shopID}, mock.Anything).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			shopPromotionRepo := mocks.NewShopPromotionRepository(t)
			tc.beforeTest(shopService, shopPromotionRepo)
			shopPromotionService := service.NewShopPromotionService(&service.ShopPromotionSConfig{
				ShopPromotionRepository: shopPromotionRepo,
				ShopService:             shopService,
			})

			err := shopPromotionService.UpdatePromotion(tc.input.userID, tc.input.promotionId, tc.input.request)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}
