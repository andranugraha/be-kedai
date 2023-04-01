package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSellerCategories(t *testing.T) {
	tests := []struct {
		name       string
		want       response.Response
		code       int
		beforeTest func(*mocks.ShopCategoryService)
	}{
		{
			name: "should return 200 when request is valid",
			want: response.Response{
				Code:    code.OK,
				Message: "success",
				Data: commonDto.PaginationResponse{
					Data:       []*dto.ShopCategory{},
					Page:       1,
					Limit:      10,
					TotalRows:  0,
					TotalPages: 0,
				},
			},
			code: http.StatusOK,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategories", 1, dto.GetSellerCategoriesRequest{
					Page:  1,
					Limit: 10,
				}).Return(&commonDto.PaginationResponse{
					Data:       []*dto.ShopCategory{},
					Page:       1,
					Limit:      10,
					TotalRows:  0,
					TotalPages: 0,
				}, nil)
			},
		},
		{
			name: "should return 500 when get shop categories failed",
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategories", 1, dto.GetSellerCategoriesRequest{
					Page:  1,
					Limit: 10,
				}).Return(nil, errs.ErrInternalServerError)
			},
		},
		{
			name: "should return 400 when shop not registered",
			want: response.Response{
				Code:    code.SHOP_NOT_REGISTERED,
				Message: errs.ErrShopNotFound.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategories", 1, dto.GetSellerCategoriesRequest{
					Page:  1,
					Limit: 10,
				}).Return(nil, errs.ErrShopNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tt.want)
			shopCategoryService := new(mocks.ShopCategoryService)
			tt.beforeTest(shopCategoryService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)
			handler := handler.New(&handler.HandlerConfig{
				ShopCategoryService: shopCategoryService,
			})

			c.Request = httptest.NewRequest("GET", "/sellers/categories", nil)
			handler.GetSellerCategories(c)

			assert.Equal(t, tt.code, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetSellerCategoryDetail(t *testing.T) {
	tests := []struct {
		name       string
		want       response.Response
		code       int
		beforeTest func(*mocks.ShopCategoryService)
	}{
		{
			name: "should return 200 when request is valid",
			want: response.Response{
				Code:    code.OK,
				Message: "success",
				Data:    &dto.ShopCategory{},
			},
			code: http.StatusOK,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategoryDetail", 1, 1).Return(&dto.ShopCategory{}, nil)
			},
		},
		{
			name: "should return 500 when get shop categories failed",
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategoryDetail", 1, 1).Return(nil, errs.ErrInternalServerError)
			},
		},
		{
			name: "should return 400 when shop not registered",
			want: response.Response{
				Code:    code.SHOP_NOT_REGISTERED,
				Message: errs.ErrShopNotFound.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategoryDetail", 1, 1).Return(nil, errs.ErrShopNotFound)
			},
		},
		{
			name: "should return 404 when category not found",
			want: response.Response{
				Code:    code.NOT_FOUND,
				Message: errs.ErrCategoryNotFound.Error(),
			},
			code: http.StatusNotFound,
			beforeTest: func(shopCategoryService *mocks.ShopCategoryService) {
				shopCategoryService.On("GetSellerCategoryDetail", 1, 1).Return(nil, errs.ErrCategoryNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tt.want)
			shopCategoryService := new(mocks.ShopCategoryService)
			tt.beforeTest(shopCategoryService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)
			c.AddParam("categoryId", "1")
			handler := handler.New(&handler.HandlerConfig{
				ShopCategoryService: shopCategoryService,
			})

			c.Request = httptest.NewRequest("GET", "/sellers/categories/1", nil)
			handler.GetSellerCategoryDetail(c)

			assert.Equal(t, tt.code, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}
