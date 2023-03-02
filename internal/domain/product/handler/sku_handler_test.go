package handler_test

import (
	"encoding/json"
	"errors"
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

func TestGetSKUByVariantIDs(t *testing.T) {
	type input struct {
		request    string
		beforeTest func(*mocks.SkuService)
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
			description: "should return error with status code 400 when no query param given",
			input: input{
				request:    "",
				beforeTest: func(ss *mocks.SkuService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "VariantID is required",
				},
			},
		},
		{
			description: "should return error with status code 422 when given invalid variant ID",
			input: input{
				request: "invalid_variant",
				beforeTest: func(ss *mocks.SkuService) {
					ss.On("GetSKUByVariantIDs", &dto.GetSKURequest{VariantID: "invalid_variant"}).Return(nil, errs.ErrInvalidVariantID)
				},
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_VARIANT,
					Message: errs.ErrInvalidVariantID.Error(),
				},
			},
		},
		{
			description: "should return error with status code 404 when prouct sku not found",
			input: input{
				request: "1000000,30000000",
				beforeTest: func(ss *mocks.SkuService) {
					ss.On("GetSKUByVariantIDs", &dto.GetSKURequest{VariantID: "1000000,30000000"}).Return(nil, errs.ErrSKUDoesNotExist)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errs.ErrSKUDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get product sku",
			input: input{
				request: "3",
				beforeTest: func(ss *mocks.SkuService) {
					ss.On("GetSKUByVariantIDs", &dto.GetSKURequest{VariantID: "3"}).Return(nil, errors.New("failed to get sku"))
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
			description: "should return sku with status code 200 when succeed fetching sku",
			input: input{
				request: "3,8,10",
				beforeTest: func(ss *mocks.SkuService) {
					ss.On("GetSKUByVariantIDs", &dto.GetSKURequest{VariantID: "3,8,10"}).Return(&model.Sku{}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &model.Sku{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			skuService := mocks.NewSkuService(t)
			tc.beforeTest(skuService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				SkuService: skuService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/products/skus?variantId=%s", tc.input.request), nil)

			h.GetSKUByVariantIDs(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}
