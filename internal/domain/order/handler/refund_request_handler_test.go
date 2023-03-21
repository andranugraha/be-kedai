package handler_test

import (
	"encoding/json"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRefundStatus(t *testing.T) {
	type input struct {
		invoiceId    int
		refundStatus *dto.RefundRequest
		beforeTest   func(*mocks.RefundRequestService)
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
			description: "should return error with status code 400 when binding fails",
			input: input{
				invoiceId:  1,
				beforeTest: func(m *mocks.RefundRequestService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "RefundStatus is required",
				},
			},
		},
		{
			description: "should return error with status code 400 when refund status is invalid",
			input: input{
				invoiceId: 1,
				refundStatus: &dto.RefundRequest{
					RefundStatus: "invalid",
				},
				beforeTest: func(m *mocks.RefundRequestService) {
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: commonErr.ErrInvalidRefundStatus.Error(),
				},
			},
		},
		{
			description: "should return error with status code 404 when invoice does not exist",
			input: input{
				invoiceId: 1,
				refundStatus: &dto.RefundRequest{
					RefundStatus: "SELLER_APPROVED",
				},
				beforeTest: func(m *mocks.RefundRequestService) {
					m.On("UpdateRefundStatus", 1, "SELLER_APPROVED").Return(commonErr.ErrRefundRequestNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.REFUND_REQUEST_NOT_FOUND,
					Message: commonErr.ErrRefundRequestNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when refund approval or rejection fails",
			input: input{
				invoiceId: 1,
				refundStatus: &dto.RefundRequest{
					RefundStatus: "SELLER_APPROVED",
				},
				beforeTest: func(m *mocks.RefundRequestService) {
					m.On("UpdateRefundStatus", 1, "SELLER_APPROVED").Return(commonErr.ErrInternalServerError)
				},
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
			description: "should return success with status code 200 when is successfully approved",
			input: input{
				invoiceId: 1,
				refundStatus: &dto.RefundRequest{
					RefundStatus: "SELLER_APPROVED",
				},
				beforeTest: func(m *mocks.RefundRequestService) {
					m.On("UpdateRefundStatus", 1, "SELLER_APPROVED").Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "refund status updated",
				},
			},
		},
		{
			description: "should return success with status code 200 when is successfully rejected",
			input: input{
				invoiceId: 1,
				refundStatus: &dto.RefundRequest{
					RefundStatus: "REJECTED",
				},
				beforeTest: func(m *mocks.RefundRequestService) {
					m.On("UpdateRefundStatus", 1, "REJECTED").Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "refund status updated",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			inputBody := test.MakeRequestBody(tc.input.refundStatus)
			refundRequestService := mocks.NewRefundRequestService(t)
			tc.beforeTest(refundRequestService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			c.AddParam("orderId", "1")
			c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/sellers/orders/{%d}/refund", tc.input.invoiceId), inputBody)
			handler := handler.New(&handler.Config{
				RefundRequestService: refundRequestService,
			})

			handler.UpdateRefundStatus(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}
