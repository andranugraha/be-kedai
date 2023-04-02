package service_test

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSellerCategories(t *testing.T) {
	var (
		userId = 1
		req    = dto.GetSellerCategoriesRequest{
			Page:  1,
			Limit: 10,
		}
	)
	tests := []struct {
		name       string
		req        dto.GetSellerCategoriesRequest
		want       *commonDto.PaginationResponse
		err        error
		beforeTest func(*mocks.ShopService, *mocks.ShopCategoryRepository)
	}{
		{
			name: "should return shop categories when request is valid",
			req:  req,
			want: &commonDto.PaginationResponse{
				Data:       []*dto.ShopCategory{},
				Page:       req.Page,
				Limit:      req.Limit,
				TotalRows:  0,
				TotalPages: 0,
			},
			err: nil,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetByShopID", 1, req).Return([]*dto.ShopCategory{}, int64(0), 0, nil)
			},
		},
		{
			name: "should return error when shop not found",
			req:  req,
			want: nil,
			err:  commonErr.ErrShopNotFound,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(nil, commonErr.ErrShopNotFound)
			},
		},
		{
			name: "should return error when get shop categories failed",
			req:  req,
			want: nil,
			err:  commonErr.ErrInternalServerError,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetByShopID", 1, req).Return(nil, int64(0), 0, commonErr.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shopService := new(mocks.ShopService)
			shopCategoryRepo := new(mocks.ShopCategoryRepository)

			test.beforeTest(shopService, shopCategoryRepo)

			shopCategoryService := service.NewShopCategoryService(&service.ShopCategorySConfig{
				ShopService:      shopService,
				ShopCategoryRepo: shopCategoryRepo,
			})

			got, err := shopCategoryService.GetSellerCategories(userId, test.req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.err, err)
		})
	}
}

func TestGetSellerCategoryDetail(t *testing.T) {
	var (
		userId     = 1
		categoryId = 1
	)

	tests := []struct {
		name       string
		want       *dto.ShopCategory
		err        error
		beforeTest func(*mocks.ShopService, *mocks.ShopCategoryRepository)
	}{
		{
			name: "should return shop category detail when request is valid",
			want: &dto.ShopCategory{
				ID: 1,
			},
			err: nil,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetByIDAndShopID", categoryId, 1).Return(&dto.ShopCategory{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "should return error when shop not found",
			want: nil,
			err:  commonErr.ErrShopNotFound,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(nil, commonErr.ErrShopNotFound)
			},
		},
		{
			name: "should return error when shop category not found",
			want: nil,
			err:  commonErr.ErrCategoryNotFound,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetByIDAndShopID", categoryId, 1).Return(nil, commonErr.ErrCategoryNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shopService := new(mocks.ShopService)
			shopCategoryRepo := new(mocks.ShopCategoryRepository)

			test.beforeTest(shopService, shopCategoryRepo)

			shopCategoryService := service.NewShopCategoryService(&service.ShopCategorySConfig{
				ShopService:      shopService,
				ShopCategoryRepo: shopCategoryRepo,
			})

			got, err := shopCategoryService.GetSellerCategoryDetail(userId, categoryId)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.err, err)
		})
	}
}

func TestCreateSellerCategory(t *testing.T) {
	var (
		userId = 1
		req    = dto.CreateSellerCategoryRequest{
			Name: "test",
		}
	)

	tests := []struct {
		name       string
		req        dto.CreateSellerCategoryRequest
		want       *dto.CreateSellerCategoryResponse
		err        error
		beforeTest func(*mocks.ShopService, *mocks.ShopCategoryRepository)
	}{
		{
			name: "should return shop category detail when request is valid",
			req:  req,
			want: &dto.CreateSellerCategoryResponse{
				ID: 0,
			},
			err: nil,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("Create", mock.Anything).Return(nil)
			},
		},
		{
			name: "should return error when shop not found",
			req:  req,
			want: nil,
			err:  commonErr.ErrShopNotFound,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(nil, commonErr.ErrShopNotFound)
			},
		},
		{
			name: "should return error when create shop category failed",
			req:  req,
			want: nil,
			err:  commonErr.ErrInternalServerError,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("Create", mock.Anything).Return(commonErr.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shopService := new(mocks.ShopService)
			shopCategoryRepo := new(mocks.ShopCategoryRepository)

			test.beforeTest(shopService, shopCategoryRepo)

			shopCategoryService := service.NewShopCategoryService(&service.ShopCategorySConfig{
				ShopService:      shopService,
				ShopCategoryRepo: shopCategoryRepo,
			})

			got, err := shopCategoryService.CreateSellerCategory(userId, test.req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.err, err)
		})
	}
}

func TestUpdateSellerCategory(t *testing.T) {
	var (
		userId     = 1
		categoryId = 1
		req        = dto.UpdateSellerCategoryRequest{}
	)

	tests := []struct {
		name       string
		req        dto.UpdateSellerCategoryRequest
		want       *dto.CreateSellerCategoryResponse
		err        error
		beforeTest func(*mocks.ShopService, *mocks.ShopCategoryRepository)
	}{
		{
			name: "should return shop category detail when request is valid",
			req:  req,
			want: &dto.CreateSellerCategoryResponse{
				ID: 1,
			},
			err: nil,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetCategoryByIDAndShopID", categoryId, 1).Return(&model.ShopCategory{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("Update", mock.Anything).Return(nil)
			},
		},
		{
			name: "should return error when shop not found",
			req:  req,
			want: nil,
			err:  commonErr.ErrShopNotFound,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(nil, commonErr.ErrShopNotFound)
			},
		},
		{
			name: "should return error when shop category not found",
			req:  req,
			want: nil,
			err:  commonErr.ErrCategoryNotFound,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetCategoryByIDAndShopID", categoryId, 1).Return(nil, commonErr.ErrCategoryNotFound)
			},
		},
		{
			name: "should return error when update shop category failed",
			req:  req,
			want: nil,
			err:  commonErr.ErrInternalServerError,
			beforeTest: func(shopService *mocks.ShopService, shopCategoryRepo *mocks.ShopCategoryRepository) {
				shopService.On("FindShopById", userId).Return(&model.Shop{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("GetCategoryByIDAndShopID", categoryId, 1).Return(&model.ShopCategory{
					ID: 1,
				}, nil)

				shopCategoryRepo.On("Update", mock.Anything).Return(commonErr.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shopService := new(mocks.ShopService)
			shopCategoryRepo := new(mocks.ShopCategoryRepository)

			test.beforeTest(shopService, shopCategoryRepo)

			shopCategoryService := service.NewShopCategoryService(&service.ShopCategorySConfig{
				ShopService:      shopService,
				ShopCategoryRepo: shopCategoryRepo,
			})

			got, err := shopCategoryService.UpdateSellerCategory(userId, categoryId, test.req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.err, err)
		})
	}
}
