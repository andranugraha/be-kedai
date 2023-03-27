package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestGetSellerPromotion(t *testing.T) {
	type input struct {
		userID   int
		request  *dto.SellerPromotionFilterRequest
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
		promotions = []*dto.SellerPromotion{}
		totalPages = 0
		request    = &dto.SellerPromotionFilterRequest{
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
			description: "should return error with status code 200 when suceed fetching promotions",
			input: input{
				userID:  userID,
				request: request,
				mockData: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       promotions,
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
						Data:       promotions,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		expectedRes, _ := json.Marshal(tc.expected.response)
		shopPromotionService := mocks.NewShopPromotionService(t)
		shopPromotionService.On("GetSellerPromotions", tc.input.userID, tc.input.request).Return(tc.input.mockData, tc.input.mockErr)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Set("userId", tc.input.userID)
		h := handler.New(&handler.HandlerConfig{
			ShopPromotionService: shopPromotionService,
		})
		c.Request = httptest.NewRequest("GET", fmt.Sprintf("/v1/sellers/promotions?page=%d&limit=%d", tc.input.request.Page, tc.input.request.Limit), nil)

		h.GetSellerPromotions(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}
}
