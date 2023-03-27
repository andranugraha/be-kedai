package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/common/dto"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	categoryDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/handler"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/server"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"

	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	var (
		req = categoryDto.GetCategoriesRequest{
			Limit: 0,
			Page:  1,
		}
		minPrice float64 = 100000
		res              = &dto.PaginationResponse{
			Data: []*model.Category{
				{
					ID:   1,
					Name: "Fashion",
					Children: []*model.Category{
						{
							ID:   2,
							Name: "Pria",
							Children: []*model.Category{
								{
									ID:       3,
									Name:     "Baju",
									MinPrice: &minPrice,
								},
							},
						},
					},
				},
			},
			Limit:      10,
			Page:       1,
			TotalRows:  1,
			TotalPages: 1,
		}
	)

	tests := []struct {
		name             string
		getCategoriesRes *dto.PaginationResponse
		getCategoriesErr error
		want             response.Response
		code             int
	}{
		{
			name:             "should return 200 when get categories success",
			getCategoriesRes: res,
			getCategoriesErr: nil,
			want: response.Response{
				Code:    code.OK,
				Message: "success",
				Data:    res,
			},
			code: http.StatusOK,
		},
		{
			name:             "should return 500 when get categories failed",
			getCategoriesRes: nil,
			getCategoriesErr: errorResponse.ErrInternalServerError,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errorResponse.ErrInternalServerError.Error(),
				Data:    nil,
			},
			code: http.StatusInternalServerError,
		},
		
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonRes, _ := json.Marshal(test.want)
			mockService := mocks.NewCategoryService(t)
			mockService.On("GetCategories", req).Return(test.getCategoriesRes, test.getCategoriesErr)
			productHandler := handler.New(&handler.Config{
				CategoryService: mockService,
			})
			cfg := &server.RouterConfig{
				ProductHandler: productHandler,
			}

			req, _ := http.NewRequest("GET", "/v1/products/categories", nil)
			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}

func TestAddCategory(t *testing.T) {
	var (
		req = categoryDto.CategoryDTO{
			Name:     "Fashion",
			ImageURL: "https://image.com",
		}
		res = &model.Category{
			ID:       0,
			Name:     "Fashion",
			ImageURL: "https://image.com",
			MinPrice: nil,
			ParentID: nil,
			Children: nil,
		}
	)

	tests := []struct {
		name           string
		addCategoryRes *model.Category
		addCategoryErr error
		want           response.Response
		code           int
	}{
		{
			name:           "should return 200 when add category success",
			addCategoryRes: res,
			addCategoryErr: nil,
			want: response.Response{
				Code:    code.OK,
				Message: "success",
			},
			code: http.StatusOK,
		},
		{
			name:           "should return 500 when add category failed",
			addCategoryRes: nil,
			addCategoryErr: errorResponse.ErrInternalServerError,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errorResponse.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
		},
	
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonRes, _ := json.Marshal(test.want)
			mockService := mocks.NewCategoryService(t)
			mockService.On("AddCategory", res).Return(test.addCategoryErr)
			productHandler := handler.New(&handler.Config{
				CategoryService: mockService,
			})
			cfg := &server.RouterConfig{
				ProductHandler: productHandler,
			}

			payload := testutil.MakeRequestBody(req)
			req, _ := http.NewRequest("POST", "/v1/products/categories", payload)
			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}
