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

func TestUserWishlistService_GetUserWishlist(t *testing.T) {
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
			description: "it should return error user not exist if user does not exist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrUserDoesNotExist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err:  errs.ErrUserDoesNotExist,
				data: nil,
			},
		},
		{
			description: "it should return error product not exist if product does not exist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrProductDoesNotExist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(nil, errs.ErrProductDoesNotExist)
				},
			},
			expected: expected{
				err:  errs.ErrProductDoesNotExist,
				data: nil,
			},
		},
		{
			description: "it should return error product not in wishlist if product is not in wishlist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrProductNotInWishlist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockWishlistRepo.On("GetUserWishlist", &model.UserWishlist{
						UserID:    user.ID,
						ProductID: product.ID,
					}).Return(nil, errs.ErrProductNotInWishlist)
				},
			},
			expected: expected{
				err:  errs.ErrProductNotInWishlist,
				data: nil,
			},
		},

		{
			description: "it should return user wishlist if success",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: nil,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
					mockWishlistRepo.On("GetUserWishlist", &model.UserWishlist{
						UserID:    user.ID,
						ProductID: product.ID,
					}).Return(wishlist, nil)
				},
			},
			expected: expected{
				err:  nil,
				data: wishlist,
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

			actualWishlist, actualErr := uc.GetUserWishlist(tc.input.data)

			assert.Equal(t, tc.expected.err, actualErr)
			assert.Equal(t, tc.expected.data, actualWishlist)
		})
	}
}

func TestAddUserWishlist(t *testing.T) {
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
					UserId:    user.ID,
					ProductId: product.ID,
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
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrProductDoesNotExist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(nil, errs.ErrProductDoesNotExist)
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
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrProductInWishlist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
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
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: nil,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
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
					UserId:    user.ID,
					ProductId: product.ID,
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
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrProductDoesNotExist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(nil, errs.ErrProductDoesNotExist)
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
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: errs.ErrProductNotInWishlist,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
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
					UserId:    user.ID,
					ProductId: product.ID,
				},
				err: nil,
				beforeTests: func(mockWishlistRepo *mocks.UserWishlistRepository, mockUserService *mocks.UserService, mockProductService *mocks.ProductService) {
					mockUserService.On("GetByID", user.ID).Return(user, nil)
					mockProductService.On("GetByID", product.ID).Return(product, nil)
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
