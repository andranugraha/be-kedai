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
