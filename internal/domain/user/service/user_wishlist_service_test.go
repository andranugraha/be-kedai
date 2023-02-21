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

func TestUserWishlistService_RemoveUserWishlist(t *testing.T) {
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
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "it should return error user not exist if user does not exist",
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
				err: errs.ErrUserDoesNotExist,
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
				err: errs.ErrProductDoesNotExist,
			},
		},
		{
			description: "it should return error product not in wishlist if user wishlist does not exist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      user.ID,
					ProductCode: product.Code,
				},
				err: errs.ErrProductNotInWishlist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByCode", product.Code).Return(product, nil)
					mockWishlistRepo.On("RemoveUserWishlist", wishlist).Return(errs.ErrProductNotInWishlist)
				},
			},
			expected: expected{
				err: errs.ErrProductNotInWishlist,
			},
		},
		{
			description: "it should return nil error and succesfully removed message if user wishlist is removed",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      user.ID,
					ProductCode: product.Code,
				},
				err: nil,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByCode", product.Code).Return(product, nil)
					mockWishlistRepo.On("RemoveUserWishlist", wishlist).Return(nil)
				},
			},
			expected: expected{
				err: nil,
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

			actualErr := uc.RemoveUserWishlist(tc.input.data)

			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}
