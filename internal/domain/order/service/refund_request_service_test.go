package service_test

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApproveRejectRefund(t *testing.T) {
	var emptyString = ""
	type input struct {
		invoiceId  int
		req        dto.RefundRequest
		beforeTest func(*mocks.RefundRequestRepository)
	}

	type expected struct {
		message string
		err     error
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error when fails to update refund status",
			input: input{
				invoiceId: 1,
				req: dto.RefundRequest{
					RefundStatus: "seller_approved",
				},
				beforeTest: func(m *mocks.RefundRequestRepository) {
					m.On("ApproveRejectRefund", 1, "seller_approved").Return(commonErr.ErrRefundRequestNotFound)
				},
			},
			expected: expected{
				message: emptyString,
				err:     commonErr.ErrRefundRequestNotFound,
			},
		},
		{
			description: "should return success when refund status updated",
			input: input{
				invoiceId: 1,
				req: dto.RefundRequest{
					RefundStatus: "seller_approved",
				},
				beforeTest: func(m *mocks.RefundRequestRepository) {
					m.On("ApproveRejectRefund", 1, "seller_approved").Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockRefundRequestRepo := new(mocks.RefundRequestRepository)

			c.input.beforeTest(mockRefundRequestRepo)
			refundRequestService := service.NewRefundRequestService(&service.RefundRequestSConfig{
				RefundRequestRepo: mockRefundRequestRepo,
			})

			err := refundRequestService.UpdateRefundStatus(c.input.invoiceId, c.input.req.RefundStatus)

			assert.Equal(t, c.expected.err, err)
		})
	}
}
