package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetInvoicesByUserID(t *testing.T) {
	type input struct {
		userID     int
		request    *dto.InvoicePerShopFilterRequest
		beforeTest func(*mocks.InvoicePerShopService)
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
			description: "should return error with status code 400 when given invalid queries",
			input: input{
				userID: 1,
				request: &dto.InvoicePerShopFilterRequest{
					StartDate: "test",
				},
				beforeTest: func(ipss *mocks.InvoicePerShopService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "EndDate is required",
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get invoices",
			input: input{
				userID:  1,
				request: &dto.InvoicePerShopFilterRequest{},
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByUserID", 1, &dto.InvoicePerShopFilterRequest{Limit: 10, Page: 1}).Return(nil, errors.New("failed to get invoices"))
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
			description: "should return error with status code 200 when suceed fetching invoices",
			input: input{
				userID:  1,
				request: &dto.InvoicePerShopFilterRequest{},
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByUserID", 1, &dto.InvoicePerShopFilterRequest{Limit: 10, Page: 1}).Return(&commonDto.PaginationResponse{Page: 1, Limit: 10, Data: []*dto.InvoicePerShopDetail{}}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &commonDto.PaginationResponse{Page: 1, Limit: 10, Data: []*dto.InvoicePerShopDetail{}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			tc.beforeTest(invoicePerShopService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/invoices?startDate=%s&endDate=%s&s=%s&status=%s", tc.input.request.StartDate, tc.input.request.EndDate, tc.input.request.S, tc.input.request.Status), nil)
			handler := handler.New(&handler.HandlerConfig{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.GetInvoicePerShopsByUserID(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}
