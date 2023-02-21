package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	model "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	mocks "kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserWishlistService_AddUserWishlist(t *testing.T) {
	var (
		product = &productModel.Product{
			ID:   1,
			Code: "123",
		}
		user = &model.User{
			ID: 1,
		}
		wishlist = &model.UserWishlist{
			UserID:    user.ID,
			ProductID: product.ID,
		}
	)
	type input struct {
		data        *dto.UserWishlistRequest
		err         error
		beforeTests func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService)
	}
	type expected struct {
		data *model.UserWishlist
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "it should return nil data and error user not exist if user does not exist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      user.ID,
					ProductCode: product.Code,
				},
				err: errs.ErrUserDoesNotExist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "it should return nil data and error product not exist if product does not exist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      user.ID,
					ProductCode: product.Code,
				},
				err: errs.ErrProductDoesNotExist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByCode", product.Code).Return(nil, errs.ErrProductDoesNotExist)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrProductDoesNotExist,
			},
		},
		{
			description: "it should return nil data and error product in wishlist if user wishlist aleardy exists",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      user.ID,
					ProductCode: product.Code,
				},
				err: errs.ErrProductInWishlist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByCode", product.Code).Return(product, nil)
					mockWishlistRepo.On("AddUserWishlist", wishlist).Return(nil, errs.ErrProductInWishlist)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrProductInWishlist,
			},
		},
		{
			description: "it should return wishlist data and nil error",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      user.ID,
					ProductCode: product.Code,
				},
				err: nil,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByCode", product.Code).Return(product, nil)
					mockWishlistRepo.On("AddUserWishlist", wishlist).Return(wishlist, nil)
				},
			},
			expected: expected{
				data: wishlist,
				err:  nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			mockWishlistRepo := mocks.NewUserWishlistRepository(t)
			mockUserService := mocks.NewUserService(t)
			mockProductService := mocks.NewProductService(t)

			tc.beforeTests(mockWishlistRepo, mockUserService, mockProductService)

			uc := service.NewUserWishlistService(&service.UserWishlistSConfig{
				UserWishlistRepository: mockWishlistRepo,
				UserService:            mockUserService,
				ProductService:         mockProductService,
			})

			actualUserWishlist, actualErr := uc.AddUserWishlist(tc.input.data)

			assert.Equal(t, tc.expected.data, actualUserWishlist)
			assert.Equal(t, actualErr, tc.expected.err)
		})
	}
}
