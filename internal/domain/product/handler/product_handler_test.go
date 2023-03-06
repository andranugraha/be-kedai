package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRecommendation(t *testing.T) {
	var (
		req = dto.RecommendationByCategoryIdRequest{
			CategoryId: 2,
			ProductId:  2,
		}
		invalidReq = dto.RecommendationByCategoryIdRequest{
			CategoryId: 1,
		}
		products = []*dto.ProductResponse{}
	)

	type input struct {
		dto     dto.RecommendationByCategoryIdRequest
		product []*dto.ProductResponse
		err     error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return recommended product list with code 200 when success",
			input: input{
				dto:     req,
				product: products,
				err:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    products,
				},
			},
		},
		{
			description: "should return error with code 400 when required param not met",
			input: input{
				dto:     invalidReq,
				product: nil,
				err:     errs.ErrBadRequest,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "ProductId is required",
				},
			},
		},
		{
			description: "should return error with code 400 when category id not exist",
			input: input{
				dto:     req,
				product: nil,
				err:     errs.ErrCategoryDoesNotExist,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "category doesn't exist",
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				dto:     req,
				product: nil,
				err:     errs.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: "something went wrong in the server",
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.ProductService)
			mockService.On("GetRecommendationByCategory", tc.input.dto.ProductId, tc.input.dto.CategoryId).Return(tc.input.product, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				ProductService: mockService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/products/recommendation?productId=%d&categoryId=%d", tc.dto.ProductId, tc.dto.CategoryId), nil)

			h.GetRecommendationByCategory(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetProductByCode(t *testing.T) {
	type input struct {
		productCode string
		mockReturn  *dto.ProductDetail
		mockErr     error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return with status code 404 if product not found",
			input: input{
				productCode: "product_code",
				mockReturn:  nil,
				mockErr:     errs.ErrProductDoesNotExist,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.PRODUCT_NOT_EXISTS,
					Message: errs.ErrProductDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return with status code 500 if product not found",
			input: input{
				productCode: "product_code",
				mockReturn:  nil,
				mockErr:     errors.New("failed to get product"),
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return with status code 200 if product found",
			input: input{
				productCode: "product_code",
				mockReturn:  &dto.ProductDetail{},
				mockErr:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    &dto.ProductDetail{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			productServiceMock := mocks.NewProductService(t)
			productServiceMock.On("GetByCode", tc.input.productCode).Return(tc.input.mockReturn, tc.input.mockErr)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.AddParam("code", tc.input.productCode)
			h := handler.New(&handler.Config{
				ProductService: productServiceMock,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/v1/products/%s", tc.input.productCode), nil)

			h.GetProductByCode(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestProductSearchFiltering(t *testing.T) {
	var (
		product = []*dto.ProductResponse{}
		req     = dto.ProductSearchFilterRequest{
			Limit: 10,
			Page:  1,
			Sort:  "recommended",
		}
		res = &commonDto.PaginationResponse{
			Data:  product,
			Limit: 10,
			Page:  1,
		}
	)
	type input struct {
		dto     dto.ProductSearchFilterRequest
		product *commonDto.PaginationResponse
		err     error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return filtered product list with code 200 when success",
			input: input{
				dto:     req,
				product: res,
				err:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    res,
				},
			},
		},
		{
			description: "should return error code 404 when shop not found",
			input: input{
				dto:     req,
				product: nil,
				err:     errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				dto:     req,
				product: nil,
				err:     errs.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: "something went wrong in the server",
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.ProductService)
			mockService.On("ProductSearchFiltering", tc.input.dto).Return(tc.input.product, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				ProductService: mockService,
			})
			c.Request = httptest.NewRequest("GET", "/products", nil)

			h.ProductSearchFiltering(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}
