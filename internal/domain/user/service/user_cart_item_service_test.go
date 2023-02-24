package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	mocks "kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreCheckCartItem(t *testing.T) {
	var (
		req = &dto.UserCartItemRequest{
			Quantity: 1,
			Notes:    "notes",
			UserId:   1,
			SkuId:    1,
		}
		sku = &productModel.Sku{
			ID:        1,
			ProductId: 1,
			Stock:     10,
		}
		product = &productModel.Product{
			ID:       1,
			ShopID:   1,
			IsActive: true,
		}
	)
	type input struct {
		data        *dto.UserCartItemRequest
		err         error
		beforeTests func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService)
	}
	type expected struct {
		data *model.CartItem
		sku  *productModel.Sku
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error if sku not found",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(nil, errs.ErrProductDoesNotExist)
				},
				err: errs.ErrProductDoesNotExist,
			},
			expected: expected{
				err: errs.ErrProductDoesNotExist,
			},
		},
		{
			description: "should return error if product not exist",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(nil, errs.ErrProductDoesNotExist)
				},
				err: errs.ErrProductDoesNotExist,
			},
			expected: expected{
				err: errs.ErrProductDoesNotExist,
			},
		},
		{
			description: "should return error if product not active",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(&productModel.Product{
						ID:       product.ID,
						IsActive: false,
					}, nil)
				},
				err: errs.ErrProductDoesNotExist,
			},
			expected: expected{
				err: errs.ErrProductDoesNotExist,
			},
		},
		{
			description: "should return error if internal server error",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(nil, errs.ErrInternalServerError)
				},
				err: errs.ErrInternalServerError,
			},
			expected: expected{
				err:  errs.ErrInternalServerError,
				data: nil,
				sku:  nil,
			},
		},
		{
			description: "should return error if shop user id is same as cart item user id",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 1,
					}, nil)
				},
				err: errs.ErrUserIsShopOwner,
			},
			expected: expected{
				err:  errs.ErrUserIsShopOwner,
				data: nil,
				sku:  nil,
			},
		},
		{
			description: "should return error other than ErrCartItemNotFound",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				err:  errs.ErrInternalServerError,
				data: nil,
				sku:  nil,
			},
		},
		{
			description: "should return error if cart item quantity is greater than sku quantity",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 10,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(&model.CartItem{
						Quantity: 1,
					}, nil)
				},
				err: errs.ErrProductQuantityNotEnough,
			},
			expected: expected{
				err:  errs.ErrProductQuantityNotEnough,
				data: nil,
				sku:  nil,
			},
		},
		{
			description: "should return cart item, sku and nil error if quantity is less than sku quantity",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(&model.CartItem{
						Quantity: 2,
					}, nil)
				},
			},
			expected: expected{
				err: nil,
				data: &model.CartItem{
					Quantity: 2,
				},
				sku: sku,
			},
		},
		{
			description: "should return nil cart item, sku and nil error if cart item not found",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(nil, errs.ErrCartItemNotFound)
				},
			},
			expected: expected{
				err:  nil,
				data: nil,
				sku:  sku,
			},
		},
		{
			description: "should return nil cart item, nil sku and error if other error happened",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, mockShopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockShopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				err:  errs.ErrInternalServerError,
				data: nil,
				sku:  nil,
			}}}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockCartItemRepo := mocks.NewUserCartItemRepository(t)
			mockSkuService := mocks.NewSkuService(t)
			mockProductService := mocks.NewProductService(t)
			mockShopService := mocks.NewShopService(t)

			c.beforeTests(mockCartItemRepo, mockSkuService, mockProductService, mockShopService)

			s := service.NewUserCartItemService(&service.UserCartItemSConfig{
				CartItemRepository: mockCartItemRepo,
				SkuService:         mockSkuService,
				ProductService:     mockProductService,
				ShopService:        mockShopService,
			})

			result, sku, err := s.PreCheckCartItem(c.input.data)

			assert.ErrorIs(t, err, c.expected.err)
			assert.Equal(t, c.expected.data, result)
			assert.Equal(t, c.expected.sku, sku)
		})
	}

}

func TestCreateCartItem(t *testing.T) {
	var (
		req = &dto.UserCartItemRequest{
			Quantity: 1,
			Notes:    "notes",
			UserId:   1,
			SkuId:    1,
		}
		sku = &productModel.Sku{
			ID:        1,
			ProductId: 1,
			Stock:     10,
		}
		product = &productModel.Product{
			ID:       1,
			ShopID:   1,
			IsActive: true,
		}
	)
	type input struct {
		data        *dto.UserCartItemRequest
		err         error
		beforeTests func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService)
	}
	type expected struct {
		data *model.CartItem
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return nil cart item and error if other error happened",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					shopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(nil, errs.ErrInternalServerError)

				},
				err: errs.ErrInternalServerError,
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return nil cart item and error if other error happened",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					shopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(nil, errs.ErrCartItemNotFound)
					mockUserCartItemRepo.On("CreateCartItem", &model.CartItem{
						Quantity: 1,
						SkuId:    1,
						UserId:   1,
					}).Return(nil, errs.ErrInternalServerError)
				},
				err: errs.ErrInternalServerError,
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return cart item and nil error if success",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					shopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(nil, errs.ErrCartItemNotFound)
					mockUserCartItemRepo.On("CreateCartItem", &model.CartItem{
						Quantity: 1,
						SkuId:    1,
						UserId:   1,
					}).Return(&model.CartItem{
						Quantity: 1,
						SkuId:    1,
						UserId:   1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.CartItem{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				err: nil,
			},
		},
		{
			description: "should return nil cart item and error if update cart item failed",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					shopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(&model.CartItem{
						Quantity: 1,
						SkuId:    1,
						UserId:   1,
					}, nil)
					mockUserCartItemRepo.On("UpdateCartItem", &model.CartItem{
						Quantity: 2,
						SkuId:    1,
						UserId:   1,
					}).Return(&model.CartItem{
						Quantity: 2,
						SkuId:    1,
						UserId:   1,
					}, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},

		{
			description: "should return cart item and nil error if update cart item success",
			input: input{
				data: &dto.UserCartItemRequest{
					Quantity: 1,
					SkuId:    1,
					UserId:   1,
					Notes:    "test",
				},
				beforeTests: func(mockUserCartItemRepo *mocks.UserCartItemRepository, mockSkuService *mocks.SkuService, mockProductService *mocks.ProductService, shopService *mocks.ShopService) {
					mockSkuService.On("GetByID", sku.ID).Return(sku, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					shopService.On("FindShopByUserId", req.UserId).Return(&shopModel.Shop{
						ID:     1,
						UserID: 2,
					}, nil)
					mockUserCartItemRepo.On("GetCartItemByUserIdAndSkuId", req.UserId, sku.ID).Return(&model.CartItem{
						Quantity: 1,
						SkuId:    1,
						UserId:   1,
					}, nil)
					mockUserCartItemRepo.On("UpdateCartItem", &model.CartItem{
						Quantity: 2,
						SkuId:    1,
						UserId:   1,
						Notes:    "test",
					}).Return(&model.CartItem{
						Quantity: 2,
						SkuId:    1,
						UserId:   1,
						Notes:    "test",
					}, nil)
				},
			},
			expected: expected{
				data: &model.CartItem{
					Quantity: 2,
					SkuId:    1,
					UserId:   1,
					Notes:    "test",
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockCartItemRepo := mocks.NewUserCartItemRepository(t)
			mockSkuService := mocks.NewSkuService(t)
			mockProductService := mocks.NewProductService(t)
			mockShopService := mocks.NewShopService(t)

			c.beforeTests(mockCartItemRepo, mockSkuService, mockProductService, mockShopService)

			s := service.NewUserCartItemService(&service.UserCartItemSConfig{
				CartItemRepository: mockCartItemRepo,
				SkuService:         mockSkuService,
				ProductService:     mockProductService,
				ShopService:        mockShopService,
			})

			result, err := s.CreateCartItem(c.input.data)

			assert.ErrorIs(t, err, c.expected.err)
			assert.Equal(t, c.expected.data, result)
		})
	}

}
