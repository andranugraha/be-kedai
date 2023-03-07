package handler_test

import (
	"encoding/json"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
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

func TestGetProductReviews(t *testing.T) {
	type input struct {
		req dto.GetReviewRequest
		res commonDto.PaginationResponse
		err error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return internal server error",
			input: input{
				req: dto.GetReviewRequest{
					Limit: 6,
					Page:  1,
				},
				res: commonDto.PaginationResponse{},
				err: commonErr.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: commonErr.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return success",
			input: input{
				req: dto.GetReviewRequest{
					Limit: 6,
					Page:  1,
				},
				res: commonDto.PaginationResponse{
					Data: []interface{}{},
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data: commonDto.PaginationResponse{
						Data: []interface{}{},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.TransactionReviewService)
			mockService.On("GetReviews", tc.input.req).Return(&tc.input.res, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				TransactionReviewService: mockService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/products/{%s}/reviews", "test"), nil)

			h.GetProductReviews(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}

}

func TestGetProductReviewStats(t *testing.T) {
	type input struct {
		code string
		res  dto.GetReviewStatsResponse
		err  error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return internal server error",
			input: input{
				code: "test",
				res:  dto.GetReviewStatsResponse{},
				err:  commonErr.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: commonErr.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return success",
			input: input{
				code: "test",
				res:  dto.GetReviewStatsResponse{},
				err:  nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    dto.GetReviewStatsResponse{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.TransactionReviewService)
			mockService.On("GetReviewStats", tc.input.code).Return(&tc.input.res, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				TransactionReviewService: mockService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/products/{%s}/reviews/stats", tc.input.code), nil)
			c.Params = gin.Params{
				{
					Key:   "code",
					Value: tc.input.code,
				},
			}

			h.GetProductReviewStats(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}

}
