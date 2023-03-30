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
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		request     *dto.UpdateProductActivationRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID      = 1
		productCode = "product-code"
		isActive    = false
		request     = &dto.UpdateProductActivationRequest{
			IsActive: &isActive,
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
				request:     &dto.UpdateProductActivationRequest{},
			},
			beforeTest: func(ps *mocks.ProductService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "IsActive is required",
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

func TestCreateProduct(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateProductRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID              = 1
		productName         = "product name"
		description         = "product description. Fill the rest here..."
		isHazardous         = false
		isActive            = false
		isNew               = true
		weight      float64 = 1
		length      float64 = 1
		height      float64 = 1
		width       float64 = 1
		categoryID          = 1
		media               = []string{"http://test.image.png"}
		courierIDs          = []int{1}
		stock               = 1
		price       float64 = 1
		request             = &dto.CreateProductRequest{
			Name:        productName,
			Description: description,
			IsHazardous: &isHazardous,
			IsActive:    &isActive,
			IsNew:       &isNew,
			Weight:      weight,
			Width:       width,
			Height:      height,
			Length:      length,
			CategoryID:  categoryID,
			Media:       media,
			CourierIDs:  courierIDs,
			Stock:       stock,
			Price:       price,
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ProductService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.CreateProductRequest{
					Name:        "a",
					Description: description,
					IsHazardous: &isHazardous,
					IsActive:    &isActive,
					IsNew:       &isNew,
					Weight:      weight,
					Width:       width,
					Height:      height,
					Length:      length,
					CategoryID:  categoryID,
					Media:       media,
					CourierIDs:  courierIDs,
					Stock:       stock,
					Price:       price,
				},
			},
			beforeTest: func(ps *mocks.ProductService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Name must be greater than 5",
				},
			},
		},
		{
			description: "should return error with status code 404 when failed to get shop",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("CreateProduct", userID, request).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with status code 409 when sku already used",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("CreateProduct", userID, request).Return(nil, errs.ErrSKUUsed)
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.SKU_USED,
					Message: errs.ErrSKUUsed.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when product name is invalid",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("CreateProduct", userID, request).Return(nil, errs.ErrInvalidProductNamePattern)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PRODUCT_NAME,
					Message: errs.ErrInvalidProductNamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to create product",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("CreateProduct", userID, request).Return(nil, errors.New("failed to create product"))
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
			description: "should return error with status code 201 when succeed to create product",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("CreateProduct", userID, request).Return(&model.Product{}, nil)
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "product created",
					Data:    &model.Product{},
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
			h := handler.New(&handler.Config{
				ProductService: productService,
			})
			payload := test.MakeRequestBody(tc.input.request)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/products", payload)

			h.CreateProduct(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestGetRecommendedProducts(t *testing.T) {
	type input struct {
		request dto.GetRecommendedProductRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}
	var (
		defaultLimit                = 18
		recommendedProductsResponse = &commonDto.PaginationResponse{
			Data:       []*dto.ProductResponse{},
			Limit:      defaultLimit,
			Page:       1,
			TotalRows:  1,
			TotalPages: 1,
		}
		request = dto.GetRecommendedProductRequest{
			Limit: defaultLimit,
			Page:  1,
		}
	)
	test := []struct {
		description string
		input
		beforeTest func(*mocks.ProductService)
		expected
	}{
		{
			description: "should return error with status code 500 when server fails to fetch recommended products",
			input: input{
				request: dto.GetRecommendedProductRequest{
					Limit: defaultLimit,
				},
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("GetRecommendedProducts", &request).Return(nil, errors.New("failed to fetch recommended products"))
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
			description: "should return success with status code 200 when server succeed to fetch recommended products",
			input: input{
				request: dto.GetRecommendedProductRequest{
					Limit: defaultLimit,
				},
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("GetRecommendedProducts", &request).Return(recommendedProductsResponse, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    recommendedProductsResponse,
				},
			},
		},
	}
	for _, tc := range test {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			productService := mocks.NewProductService(t)
			tc.beforeTest(productService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.AddParam("limit", strconv.Itoa(tc.input.request.Limit))
			c.AddParam("page", strconv.Itoa(tc.input.request.Page))
			c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/products/recommended?limit=%d&page=%d", tc.input.request.Limit, tc.input.request.Page), nil)
			h := handler.New(&handler.Config{
				ProductService: productService,
			})
			h.GetRecommendedProducts(c)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	var (
		userID      = 1
		productCode = "product-code"
		isHazardous = false
		isNew       = true
		isActive    = true
		bulkPrice   = dto.ProductBulkPriceRequest{
			MinQuantity: 10,
			Price:       10000,
		}
		invalidMedia   = []string{"media1", "media2"}
		validMedia     = []string{"https://image.png", "https://image2.png"}
		courierIDs     = []int{1, 2, 3}
		variantGroups  = []*dto.CreateVariantGroupRequest{}
		invalidRequest = dto.CreateProductRequest{
			Name:        "product-name",
			Description: "product-description",
		}
		skus              = []*dto.CreateSKURequest{}
		invalidUrlRequest = dto.CreateProductRequest{
			Name:          "product-name",
			Description:   "product-description here is more than 20 words",
			Price:         10000,
			IsHazardous:   &isHazardous,
			Weight:        1000,
			Length:        1000,
			Width:         1000,
			Height:        1000,
			IsNew:         &isNew,
			IsActive:      &isActive,
			CategoryID:    1,
			BulkPrice:     &bulkPrice,
			Media:         invalidMedia,
			CourierIDs:    courierIDs,
			Stock:         100,
			VariantGroups: variantGroups,
			SKU:           skus,
		}
		invalidNameRequest = dto.CreateProductRequest{
			Name:          "127.0.0.1",
			Description:   "product-description here is more than 20 words",
			Price:         10000,
			IsHazardous:   &isHazardous,
			Weight:        1000,
			Length:        1000,
			Width:         1000,
			Height:        1000,
			IsNew:         &isNew,
			IsActive:      &isActive,
			CategoryID:    1,
			BulkPrice:     &bulkPrice,
			Media:         validMedia,
			CourierIDs:    courierIDs,
			Stock:         100,
			VariantGroups: variantGroups,
			SKU:           skus,
		}
		request = dto.CreateProductRequest{
			Name:          "product-name",
			Description:   "product-description here is more than 20 words",
			Price:         10000,
			IsHazardous:   &isHazardous,
			Weight:        1000,
			Length:        1000,
			Width:         1000,
			Height:        1000,
			IsNew:         &isNew,
			IsActive:      &isActive,
			CategoryID:    1,
			BulkPrice:     &bulkPrice,
			Media:         validMedia,
			CourierIDs:    courierIDs,
			Stock:         100,
			VariantGroups: variantGroups,
			SKU:           skus,
		}
		updatedProduct = &model.Product{
			ID:          1,
			Code:        productCode,
			Name:        "product-name",
			Description: "product-description here is more than 20 words",
		}
	)
	type input struct {
		userID  int
		code    string
		request *dto.CreateProductRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}
	cases := []struct {
		description string
		input
		beforeTest func(*mocks.ProductService)
		expected
	}{
		{
			description: "should return error with status code 400 when request binding fails",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &invalidRequest,
			},
			beforeTest: func(ps *mocks.ProductService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Price is required",
				},
			},
		},
		{
			description: "should return error with status code 400 when product url is not valid",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &invalidUrlRequest,
			},
			beforeTest: func(ps *mocks.ProductService) {
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Media[1] must be a URL",
				},
			},
		},
		{
			description: "should return error with status code 404 when product does not exist in user's shop",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProduct", userID, productCode, &request).Return(nil, errs.ErrProductDoesNotExist)
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
			description: "should return error with status code 404 when user does not own a shop",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProduct", userID, productCode, &request).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with status code 422 when product name is invalid",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &invalidNameRequest,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProduct", userID, productCode, &invalidNameRequest).Return(nil, errs.ErrInvalidProductNamePattern)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PRODUCT_NAME,
					Message: errs.ErrInvalidProductNamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when server errors",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProduct", userID, productCode, &request).Return(nil, errs.ErrInternalServerError)
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
			description: "should return updated product with status code 200 on success",
			input: input{
				userID:  userID,
				code:    productCode,
				request: &request,
			},
			beforeTest: func(ps *mocks.ProductService) {
				ps.On("UpdateProduct", userID, productCode, &request).Return(updatedProduct, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "product updated",
					Data:    updatedProduct,
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {

			payload := test.MakeRequestBody(tc.input.request)
			expectedRes, _ := json.Marshal(tc.expected.response)
			productService := mocks.NewProductService(t)
			tc.beforeTest(productService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("code", tc.input.code)
			c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/products/%s", tc.input.code), payload)
			h := handler.New(&handler.Config{
				ProductService: productService,
			})
			h.UpdateProduct(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}
