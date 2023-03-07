package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/handler"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFindShopBySlug(t *testing.T) {
	var (
		slug       = "shop"
		shopResult = &model.Shop{
			ID: 1,
		}
	)

	type input struct {
		shop *model.Shop
		err  error
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
			description: "should return shop information with code 200 when success",
			input: input{
				shop: shopResult,
				err:  nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    shopResult,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				shop: nil,
				err:  errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: "shop not found",
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				shop: nil,
				err:  errs.ErrInternalServerError,
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
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{
				{
					Key:   "slug",
					Value: slug,
				},
			}
			mockService := new(mocks.ShopService)
			mockService.On("FindShopBySlug", slug).Return(tc.input.shop, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/shops/:slug", nil)

			handler.FindShopBySlug(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestFindShopByKeyword(t *testing.T) {
	var (
		shopList   = []*model.Shop{}
		pagination = &commonDto.PaginationResponse{
			Data:  shopList,
			Limit: 10,
			Page:  1,
		}
		req = dto.FindShopRequest{
			Keyword: "test",
			Page:    1,
			Limit:   10,
		}
	)
	type input struct {
		dto    dto.FindShopRequest
		result *commonDto.PaginationResponse
		err    error
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
			description: "should return shop list with code 200 when success",
			input: input{
				dto:    req,
				result: pagination,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    pagination,
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				dto:    req,
				result: nil,
				err:    errs.ErrInternalServerError,
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
			expectedBody, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockService := new(mocks.ShopService)
			mockService.On("FindShopByKeyword", tc.input.dto).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/shops?keyword=test", nil)

			handler.FindShopByKeyword(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}
