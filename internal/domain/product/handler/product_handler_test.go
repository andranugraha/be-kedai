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
	"kedai/backend/be-kedai/internal/utils/test"
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

func TestGetProductsByShopSlug(t *testing.T) {
	type input struct {
		slug       string
		request    *dto.ShopProductFilterRequest
		mockReturn *commonDto.PaginationResponse
		mockErr    error
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
			description: "should return error with status code 404 when shop not found",
			input: input{
				slug: "invalid-slug",
				request: &dto.ShopProductFilterRequest{
					Sort:      "recommended",
					PriceSort: "price_low",
					Limit:     10,
					Page:      1,
				},
				mockReturn: nil,
				mockErr:    errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get shop product list",
			input: input{
				slug: "invalid-slug",
				request: &dto.ShopProductFilterRequest{
					Sort:      "recommended",
					PriceSort: "price_low",
					Limit:     10,
					Page:      1,
				},
				mockReturn: nil,
				mockErr:    errors.New("failed to get shop product list"),
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
			description: "should return shop product list data with status code 200 when fetching data succeed",
			input: input{
				slug: "invalid-slug",
				request: &dto.ShopProductFilterRequest{
					Sort:      "recommended",
					PriceSort: "price_low",
					Limit:     10,
					Page:      1,
				},
				mockReturn: &commonDto.PaginationResponse{},
				mockErr:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &commonDto.PaginationResponse{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.ProductService)
			mockService.On("GetProductsByShopSlug", tc.input.slug, tc.input.request).Return(tc.input.mockReturn, tc.input.mockErr)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.AddParam("slug", tc.input.slug)
			h := handler.New(&handler.Config{
				ProductService: mockService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/shops/%s/products", tc.input.slug), nil)

			h.GetProductsByShopSlug(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestSearchAutocomplete(t *testing.T) {
	var (
		req = dto.ProductSearchAutocomplete{
			Limit: 10,
		}
		res = []*dto.ProductResponse{}
	)
	type input struct {
		req dto.ProductSearchAutocomplete
		res []*dto.ProductResponse
		err error
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
			description: "should return list of matched product with code 200 when success",
			input: input{
				req: req,
				res: res,
				err: nil,
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
			description: "should return error with code 500 when internal server error",
			input: input{
				req: req,
				res: nil,
				err: errs.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.ProductService)
			mockService.On("SearchAutocomplete", tc.input.req).Return(tc.input.res, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				ProductService: mockService,
			})
			c.Request = httptest.NewRequest("GET", "/products/autocompletes", nil)

			h.SearchAutocomplete(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestGetSellerProduct(t *testing.T) {
	type input struct {
		userID   int
		request  *dto.SellerProductFilterRequest
		mockData *commonDto.PaginationResponse
		mockErr  error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID     = 1
		page       = 2
		limit      = 10
		totalRows  = int64(0)
		products   = []*dto.SellerProduct{}
		totalPages = 0
		request    = &dto.SellerProductFilterRequest{
			Page:  page,
			Limit: limit,
		}
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 404 when shop does not exist",
			input: input{
				userID:   userID,
				request:  request,
				mockData: nil,
				mockErr:  errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:   userID,
				request:  request,
				mockData: nil,
				mockErr:  errors.New("something went wrong"),
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
			description: "should return error with status code 200 when suceed fetching products",
			input: input{
				userID:  userID,
				request: request,
				mockData: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       products,
				},
				mockErr: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data: &commonDto.PaginationResponse{
						TotalRows:  totalRows,
						TotalPages: totalPages,
						Page:       page,
						Limit:      limit,
						Data:       products,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		expectedRes, _ := json.Marshal(tc.expected.response)
		productService := mocks.NewProductService(t)
		productService.On("GetSellerProducts", tc.input.userID, tc.input.request).Return(tc.input.mockData, tc.input.mockErr)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Set("userId", tc.input.userID)
		h := handler.New(&handler.Config{
			ProductService: productService,
		})
		c.Request = httptest.NewRequest("GET", fmt.Sprintf("/v1/sellers/products?page=%d&limit=%d", tc.input.request.Page, tc.input.request.Limit), nil)

		h.GetSellerProducts(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}
}

func TestGetSellerProductDetailByCode(t *testing.T) {
	type input struct {
		userID      int
		productCode string
		mockData    *dto.SellerProductDetail
		mockErr     error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID      = 1
		productCode = "product-code"
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 404 when shop not found",
			input: input{
				userID:      userID,
				productCode: productCode,
				mockData:    nil,
				mockErr:     errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 404 when product not found",
			input: input{
				userID:      userID,
				productCode: productCode,
				mockData:    nil,
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
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:      userID,
				productCode: productCode,
				mockData:    nil,
				mockErr:     errs.ErrInternalServerError,
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
			description: "should return product detail with status code 200 when succeed to get product",
			input: input{
				userID:      userID,
				productCode: productCode,
				mockData:    &dto.SellerProductDetail{},
				mockErr:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &dto.SellerProductDetail{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			productService := mocks.NewProductService(t)
			productService.On("GetSellerProductByCode", tc.input.userID, tc.input.productCode).Return(tc.input.mockData, tc.input.mockErr)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("code", tc.input.productCode)
			h := handler.New(&handler.Config{
				ProductService: productService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/sellers/products?%s", tc.input.productCode), nil)

			h.GetSellerProductDetailByCode(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestAddProductView(t *testing.T) {
	type input struct {
		req        dto.AddProductViewRequest
		beforeTest func(mockProductService *mocks.ProductService)
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		productID = 1
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 400 when request body is invalid",
			input: input{
				req: dto.AddProductViewRequest{},
				beforeTest: func(mockProductService *mocks.ProductService) {
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "ProductID is required",
				},
			},
		},
		{
			description: "should return error with status code 500 when something went wrong",
			input: input{
				req: dto.AddProductViewRequest{
					ProductID: productID,
				},
				beforeTest: func(mockProductService *mocks.ProductService) {
					mockProductService.On("AddViewCount", productID).Return(errs.ErrInternalServerError)
				},
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
			description: "should return error with status code 404 when product not found",
			input: input{
				req: dto.AddProductViewRequest{
					ProductID: productID,
				},
				beforeTest: func(mockProductService *mocks.ProductService) {
					mockProductService.On("AddViewCount", productID).Return(errs.ErrProductDoesNotExist)
				},
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
			description: "should return status code 200 when suceed adding product view",
			input: input{
				req: dto.AddProductViewRequest{
					ProductID: productID,
				},

				beforeTest: func(mockProductService *mocks.ProductService) {
					mockProductService.On("AddViewCount", productID).Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
	}

	for _, tc := range tests {
		expectedRes, _ := json.Marshal(tc.expected.response)
		productService := mocks.NewProductService(t)
		tc.beforeTest(productService)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		payload := test.MakeRequestBody(tc.input.req)
		c.Request = httptest.NewRequest("POST", "/v1/products/views", payload)
		h := handler.New(&handler.Config{
			ProductService: productService,
		})

		h.AddProductView(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}

}

func TestXxx(t *testing.T) {
	type input struct {
		userID      int
		productCode string
		request     *dto.UpdateProductActiationRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID      = 1
		productCode = "product-code"
		request     = &dto.UpdateProductActiationRequest{
			IsActive: "true",
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ProductService)
		expected
	}{
		{
			description: "should return error with status code 400 when given empty request body",
			input: input{
				userID:      userID,
				productCode: productCode,
				request: &dto.UpdateProductActiationRequest{
					IsActive: "not-boolean",
				},
			},
			beforeTest: func(ps *mocks.ProductService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "IsActive must be either 'true' or 'false'",
				},
			},
		},
		{
			description: "should return error with status code 404 when shop not found",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProductActivation", userID, productCode, request).Return(errs.ErrShopNotFound)
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 404 when product not found",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProductActivation", userID, productCode, request).Return(errs.ErrProductDoesNotExist)
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
			description: "should return error with status code 500 when failed to update product",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProductActivation", userID, productCode, request).Return(errors.New("failed to update product activation"))
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
			description: "should return error with status code 200 when update succeed",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProductActivation", userID, productCode, request).Return(nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "update successful",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			productService := mocks.NewProductService(t)
			tc.beforeTest(productService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("code", productCode)
			h := handler.New(&handler.Config{
				ProductService: productService,
			})
			payload := test.MakeRequestBody(tc.input.request)
			c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/sellers/products/%s/activations", tc.input.productCode), payload)

			h.UpdateProductActivation(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}
