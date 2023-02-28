package handler_test

import (
	"encoding/json"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/handler"
	"kedai/backend/be-kedai/internal/domain/product/model"
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
		req = dto.RecommendationRequest{
			CategoryId: 2,
			ProductId:  2,
		}
		invalidReq = dto.RecommendationRequest{
			CategoryId: 1,
		}
		products = []*model.Product{}
	)

	type input struct {
		dto     dto.RecommendationRequest
		product []*model.Product
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
			mockService.On("GetRecommendation", tc.input.dto.ProductId, tc.input.dto.CategoryId).Return(tc.input.product, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				ProductService: mockService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/products/recommendation?productId=%d&categoryId=%d", tc.dto.ProductId, tc.dto.CategoryId), nil)

			h.GetRecommendation(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}
